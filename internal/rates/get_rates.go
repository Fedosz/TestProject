package rates

import (
	"context"
	"time"

	"rates_project/internal/domain/models"
	"rates_project/internal/domain/types"
	"rates_project/internal/telemetry"
)

func (s *Service) GetRates(
	ctx context.Context,
	askParams models.CalculationParams,
	bidParams models.CalculationParams,
) (*models.Rate, error) {
	ctx, span := telemetry.Tracer("rates_service").Start(ctx, "rates.GetRates")
	defer span.End()

	asks, bids, err := s.client.GetRates(ctx)
	if err != nil {
		return nil, err
	}

	ask, err := s.calculate(asks, askParams)
	if err != nil {
		return nil, err
	}

	bid, err := s.calculate(bids, bidParams)
	if err != nil {
		return nil, err
	}

	rate := &models.Rate{
		Ask:        ask,
		Bid:        bid,
		ReceivedAt: time.Now().UTC(),
	}

	if err = s.repo.Create(ctx, rate); err != nil {
		return nil, err
	}

	return rate, nil
}

func (s *Service) calculate(values []float64, params models.CalculationParams) (float64, error) {
	switch params.Method {
	case types.CalcMethodTopN:
		return s.calcTopN(values, params.N)
	case types.CalcMethodAvgNM:
		return s.calcAvgNM(values, params.N, params.M)
	default:
		return 0, ErrInvalidCalcMethod
	}
}
