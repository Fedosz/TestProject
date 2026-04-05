package rates

import (
	"context"
	"rates_project/internal/domain/models"
	"rates_project/internal/domain/types"
	"time"
)

func (s *Service) GetRates(
	ctx context.Context,
	askParams models.CalculationParams,
	bidParams models.CalculationParams,
) (*models.Rate, error) {
	values, err := s.client.GetRates(ctx)
	if err != nil {
		return nil, err
	}

	ask, err := s.calculate(values, askParams)
	if err != nil {
		return nil, err
	}

	bid, err := s.calculate(values, bidParams)
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
