package rates

import (
	"context"
	"rates_project/internal/domain/models"
)

type ExchangeClient interface {
	GetRates(ctx context.Context) ([]float64, error)
}

type RateRepository interface {
	Create(ctx context.Context, rate *models.Rate) error
}

type Service struct {
	client ExchangeClient
	repo   RateRepository
}

func NewService(client ExchangeClient, repo RateRepository) *Service {
	return &Service{
		client: client,
		repo:   repo,
	}
}
