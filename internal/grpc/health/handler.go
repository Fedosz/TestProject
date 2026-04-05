package health

import (
	"context"

	ratesv1 "rates_project/usdt-rates/gen/proto/rates/v1"
)

type HealthService interface {
	Check(ctx context.Context) error
}

type Handler struct {
	ratesv1.UnimplementedHealthServiceServer
	service HealthService
}

func NewHandler(service HealthService) *Handler {
	return &Handler{
		service: service,
	}
}
