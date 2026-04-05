package rates

import (
	"context"

	"rates_project/internal/domain/models"
	"rates_project/internal/metrics"
)

type ExchangeClient interface {
	GetRates(ctx context.Context) ([]float64, []float64, error)
}

type RateRepository interface {
	Create(ctx context.Context, rate *models.Rate) error
}

type Service struct {
	client  ExchangeClient
	repo    RateRepository
	metrics *metrics.Metrics
}

func NewService(
	client ExchangeClient,
	repo RateRepository,
	metrics *metrics.Metrics,
) *Service {
	return &Service{
		client:  client,
		repo:    repo,
		metrics: metrics,
	}
}
