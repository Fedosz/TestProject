package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	appHealth "rates_project/internal/app/health"
	"rates_project/internal/client/grinex"
	"rates_project/internal/config"
	dbRate "rates_project/internal/db/rate"
	grpcHealth "rates_project/internal/grpc/health"
	grpcRates "rates_project/internal/grpc/rates"
	"rates_project/internal/rates"
	ratesv1 "rates_project/usdt-rates/gen/proto/rates/v1"
)

func main() {
	cfg := config.MustLoad()

	db, err := newPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	defer func() {
		_ = db.Close()
	}()

	rateRepo := dbRate.NewRepo(db)

	grinexClient := grinex.NewClient(
		cfg.Grinex.BaseURL,
		cfg.Grinex.Path,
		cfg.Grinex.Timeout,
	)

	ratesService := rates.NewService(grinexClient, rateRepo)
	healthService := appHealth.NewService(rateRepo)

	ratesHandler := grpcRates.NewHandler(ratesService)
	healthHandler := grpcHealth.NewHandler(healthService)

	grpcServer := grpc.NewServer()

	ratesv1.RegisterRatesServiceServer(grpcServer, ratesHandler)
	ratesv1.RegisterHealthServiceServer(grpcServer, healthHandler)

	address := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Printf("gRPC server started on %s", address)

		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutdown signal received")

	stopped := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		log.Println("gRPC graceful stop timeout exceeded")
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
