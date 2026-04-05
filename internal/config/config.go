package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	GRPC     GRPCConfig
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Grinex   GrinexConfig
}

type GRPCConfig struct {
	Host string
	Port int
}

type HTTPConfig struct {
	Host string
	Port int
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type GrinexConfig struct {
	BaseURL string
	Path    string
	Timeout time.Duration
}

// MustLoad init config
func MustLoad() *Config {
	grpcHost := flag.String("grpc-host", getEnv("GRPC_HOST", "0.0.0.0"), "gRPC server host")
	grpcPort := flag.Int("grpc-port", getEnvAsInt("GRPC_PORT", 50051), "gRPC server port")

	httpHost := flag.String("http-host", getEnv("HTTP_HOST", "0.0.0.0"), "HTTP server host")
	httpPort := flag.Int("http-port", getEnvAsInt("HTTP_PORT", 2112), "HTTP server port")

	postgresHost := flag.String("postgres-host", getEnv("POSTGRES_HOST", "localhost"), "PostgreSQL host")
	postgresPort := flag.Int("postgres-port", getEnvAsInt("POSTGRES_PORT", 5432), "PostgreSQL port")
	postgresUser := flag.String("postgres-user", getEnv("POSTGRES_USER", "postgres"), "PostgreSQL user")
	postgresPassword := flag.String("postgres-password", getEnv("POSTGRES_PASSWORD", "postgres"), "PostgreSQL password")
	postgresDB := flag.String("postgres-db", getEnv("POSTGRES_DB", "rates"), "PostgreSQL database")
	postgresSSLMode := flag.String("postgres-sslmode", getEnv("POSTGRES_SSLMODE", "disable"), "PostgreSQL sslmode")

	grinexBaseURL := flag.String("grinex-base-url", getEnv("GRINEX_BASE_URL", "https://grinex.io"), "Grinex base url")
	grinexPath := flag.String("grinex-path", getEnv("GRINEX_PATH", "/api/v1/spot/depth?symbol=usdta7a5"), "Grinex API path")
	grinexTimeout := flag.Duration("grinex-timeout", getEnvAsDuration("GRINEX_TIMEOUT", 5*time.Second), "Grinex timeout")

	flag.Parse()

	return &Config{
		GRPC: GRPCConfig{
			Host: *grpcHost,
			Port: *grpcPort,
		},
		HTTP: HTTPConfig{
			Host: *httpHost,
			Port: *httpPort,
		},
		Postgres: PostgresConfig{
			Host:     *postgresHost,
			Port:     *postgresPort,
			User:     *postgresUser,
			Password: *postgresPassword,
			Database: *postgresDB,
			SSLMode:  *postgresSSLMode,
		},
		Grinex: GrinexConfig{
			BaseURL: *grinexBaseURL,
			Path:    *grinexPath,
			Timeout: *grinexTimeout,
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
