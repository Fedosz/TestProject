package rates

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"

	"rates_project/internal/domain/models"
	"rates_project/internal/domain/types"
	"rates_project/internal/rates/mocks"
)

func TestService_GetRates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		askParams models.CalculationParams
		bidParams models.CalculationParams
		prepare   func(
			exchangeClient *mocks.ExchangeClientMock,
			rateRepository *mocks.RateRepositoryMock,
		)
		check func(t *testing.T, rate *models.Rate, err error)
	}{
		{
			name: "success with top n",
			askParams: models.CalculationParams{
				Method: types.CalcMethodTopN,
				N:      1,
			},
			bidParams: models.CalculationParams{
				Method: types.CalcMethodTopN,
				N:      2,
			},
			prepare: func(
				exchangeClient *mocks.ExchangeClientMock,
				rateRepository *mocks.RateRepositoryMock,
			) {
				exchangeClient.GetRatesMock.Return(
					[]float64{80.82, 80.83, 80.84},
					[]float64{80.71, 80.70, 80.69},
					nil,
				)

				rateRepository.CreateMock.Set(func(ctx context.Context, rate *models.Rate) error {
					if rate.Ask != 80.82 {
						t.Fatalf("expected ask 80.82, got %v", rate.Ask)
					}

					if rate.Bid != 80.70 {
						t.Fatalf("expected bid 80.70, got %v", rate.Bid)
					}

					if rate.ReceivedAt.IsZero() {
						t.Fatal("expected received_at to be set")
					}

					return nil
				})
			},
			check: func(t *testing.T, rate *models.Rate, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if rate == nil {
					t.Fatal("expected rate, got nil")
				}

				if rate.Ask != 80.82 {
					t.Fatalf("expected ask 80.82, got %v", rate.Ask)
				}

				if rate.Bid != 80.70 {
					t.Fatalf("expected bid 80.70, got %v", rate.Bid)
				}

				if rate.ReceivedAt.IsZero() {
					t.Fatal("expected received_at to be set")
				}
			},
		},
		{
			name: "success with avg nm",
			askParams: models.CalculationParams{
				Method: types.CalcMethodAvgNM,
				N:      1,
				M:      2,
			},
			bidParams: models.CalculationParams{
				Method: types.CalcMethodAvgNM,
				N:      2,
				M:      3,
			},
			prepare: func(
				exchangeClient *mocks.ExchangeClientMock,
				rateRepository *mocks.RateRepositoryMock,
			) {
				exchangeClient.GetRatesMock.Return(
					[]float64{80.82, 80.83, 80.84},
					[]float64{80.71, 80.70, 80.69},
					nil,
				)

				rateRepository.CreateMock.Return(nil)
			},
			check: func(t *testing.T, rate *models.Rate, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if rate == nil {
					t.Fatal("expected rate, got nil")
				}

				assertFloatEquals(t, 80.825, rate.Ask)
				assertFloatEquals(t, 80.695, rate.Bid)
			},
		},
		{
			name: "exchange client error",
			askParams: models.CalculationParams{
				Method: types.CalcMethodTopN,
				N:      1,
			},
			bidParams: models.CalculationParams{
				Method: types.CalcMethodTopN,
				N:      1,
			},
			prepare: func(
				exchangeClient *mocks.ExchangeClientMock,
				rateRepository *mocks.RateRepositoryMock,
			) {
				exchangeClient.GetRatesMock.Return(nil, nil, errors.New("exchange error"))
			},
			check: func(t *testing.T, rate *models.Rate, err error) {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if err.Error() != "exchange error" {
					t.Fatalf("expected exchange error, got %v", err)
				}

				if rate != nil {
					t.Fatalf("expected nil rate, got %+v", rate)
				}
			},
		},
		{
			name: "repository error",
			askParams: models.CalculationParams{
				Method: types.CalcMethodTopN,
				N:      1,
			},
			bidParams: models.CalculationParams{
				Method: types.CalcMethodTopN,
				N:      1,
			},
			prepare: func(
				exchangeClient *mocks.ExchangeClientMock,
				rateRepository *mocks.RateRepositoryMock,
			) {
				exchangeClient.GetRatesMock.Return(
					[]float64{80.82, 80.83},
					[]float64{80.71, 80.70},
					nil,
				)

				rateRepository.CreateMock.Return(errors.New("repository error"))
			},
			check: func(t *testing.T, rate *models.Rate, err error) {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if err.Error() != "repository error" {
					t.Fatalf("expected repository error, got %v", err)
				}

				if rate != nil {
					t.Fatalf("expected nil rate, got %+v", rate)
				}
			},
		},
		{
			name: "invalid calculation method",
			askParams: models.CalculationParams{
				Method: "",
				N:      1,
			},
			bidParams: models.CalculationParams{
				Method: types.CalcMethodTopN,
				N:      1,
			},
			prepare: func(
				exchangeClient *mocks.ExchangeClientMock,
				rateRepository *mocks.RateRepositoryMock,
			) {
				exchangeClient.GetRatesMock.Return(
					[]float64{80.82, 80.83},
					[]float64{80.71, 80.70},
					nil,
				)
			},
			check: func(t *testing.T, rate *models.Rate, err error) {
				if err != ErrInvalidCalcMethod {
					t.Fatalf("expected %v, got %v", ErrInvalidCalcMethod, err)
				}

				if rate != nil {
					t.Fatalf("expected nil rate, got %+v", rate)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)

			exchangeClient := mocks.NewExchangeClientMock(mc)
			rateRepository := mocks.NewRateRepositoryMock(mc)

			if tt.prepare != nil {
				tt.prepare(exchangeClient, rateRepository)
			}

			service := NewService(exchangeClient, rateRepository, nil)

			rate, err := service.GetRates(
				context.Background(),
				tt.askParams,
				tt.bidParams,
			)

			tt.check(t, rate, err)
		})
	}
}

func TestService_GetRates_SetsCurrentTimestamp(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	exchangeClient := mocks.NewExchangeClientMock(mc)
	rateRepository := mocks.NewRateRepositoryMock(mc)

	exchangeClient.GetRatesMock.Return(
		[]float64{80.82},
		[]float64{80.71},
		nil,
	)

	before := time.Now().UTC()

	rateRepository.CreateMock.Set(func(ctx context.Context, rate *models.Rate) error {
		after := time.Now().UTC()

		if rate.ReceivedAt.Before(before) || rate.ReceivedAt.After(after) {
			t.Fatalf("received_at is out of expected range: %v", rate.ReceivedAt)
		}

		return nil
	})

	service := NewService(exchangeClient, rateRepository, nil)

	rate, err := service.GetRates(
		context.Background(),
		models.CalculationParams{
			Method: types.CalcMethodTopN,
			N:      1,
		},
		models.CalculationParams{
			Method: types.CalcMethodTopN,
			N:      1,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rate == nil {
		t.Fatal("expected rate, got nil")
	}

	if rate.ReceivedAt.IsZero() {
		t.Fatal("expected received_at to be set")
	}
}

func assertFloatEquals(t *testing.T, expected, actual float64) {
	t.Helper()

	const delta = 0.000001

	diff := expected - actual
	if diff < 0 {
		diff = -diff
	}

	if diff > delta {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
