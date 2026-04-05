package rates

import (
	"context"

	"rates_project/internal/domain/models"
	ratesv1 "rates_project/usdt-rates/gen/proto/rates/v1"
)

type RatesService interface {
	GetRates(
		ctx context.Context,
		askParams models.CalculationParams,
		bidParams models.CalculationParams,
	) (*models.Rate, error)
}

type Handler struct {
	ratesv1.UnimplementedRatesServiceServer
	service RatesService
}

func NewHandler(service RatesService) *Handler {
	return &Handler{
		service: service,
	}
}
