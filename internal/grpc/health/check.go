package health

import (
	"context"

	ratesv1 "rates_project/usdt-rates/gen/proto/rates/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Check(
	ctx context.Context,
	_ *ratesv1.HealthCheckRequest,
) (*ratesv1.HealthCheckResponse, error) {
	err := h.service.Check(ctx)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &ratesv1.HealthCheckResponse{
		Status: "ok",
	}, nil
}
