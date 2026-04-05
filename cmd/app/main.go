package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	"rates_project/internal/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	appHealth "rates_project/internal/app/health"
	"rates_project/internal/client/grinex"
	"rates_project/internal/config"
	"rates_project/internal/db"
	dbRate "rates_project/internal/db/rate"
	grpcHealth "rates_project/internal/grpc/health"
	grpcRates "rates_project/internal/grpc/rates"
	"rates_project/internal/rates"
	"rates_project/internal/telemetry"
	ratesv1 "rates_project/usdt-rates/gen/proto/rates/v1"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()

	log := logger.MustNew()
	defer func() {
		_ = log.Sync()
	}()

	tel, err := telemetry.Init(ctx, "rates_project")
	if err != nil {
		log.Fatal("failed to init telemetry: %v", zap.Error(err))
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = tel.Shutdown(shutdownCtx); err != nil {
			log.Warn("failed to shutdown telemetry", zap.Error(err))
		}
	}()

	postgresDB, err := newPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to connect to postgres", zap.Error(err))
	}

	defer func() {
		_ = postgresDB.Close()
	}()

	if err = db.Migrate(postgresDB, "migrations"); err != nil {
		log.Fatal("failed to run migrations", zap.Error(err))
	}

	rateRepo := dbRate.NewRepo(postgresDB)

	grinexClient := grinex.NewClient(
		cfg.Grinex.BaseURL,
		cfg.Grinex.Path,
		cfg.Grinex.Timeout,
	)

	ratesService := rates.NewService(grinexClient, rateRepo)
	healthService := appHealth.NewService(rateRepo)

	ratesHandler := grpcRates.NewHandler(ratesService)
	healthHandler := grpcHealth.NewHandler(healthService)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(telemetry.ServerStatsHandler()),
	)

	ratesv1.RegisterRatesServiceServer(grpcServer, ratesHandler)
	ratesv1.RegisterHealthServiceServer(grpcServer, healthHandler)

	address := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	stopCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Info("gRPC server started", zap.String("address", address))

		if err = grpcServer.Serve(listener); err != nil {
			log.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	<-stopCtx.Done()

	log.Info("shutdown signal received")

	stopped := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Info("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		log.Warn("gRPC graceful stop timeout exceeded")
		grpcServer.Stop()
	}
}

func newPostgresDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
