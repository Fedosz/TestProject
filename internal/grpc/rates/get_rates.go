package rates

import (
	"context"

	"rates_project/internal/domain/models"
	"rates_project/internal/domain/types"
	ratesService "rates_project/internal/rates"
	ratesv1 "rates_project/usdt-rates/gen/proto/rates/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetRates(
	ctx context.Context,
	req *ratesv1.GetRatesRequest,
) (*ratesv1.GetRatesResponse, error) {
	askParams := mapCalculationParams(req.GetAsk())
	bidParams := mapCalculationParams(req.GetBid())

	rate, err := h.service.GetRates(ctx, askParams, bidParams)
	if err != nil {
		switch err {
		case ratesService.ErrInvalidCalcMethod,
			ratesService.ErrInvalidTopN,
			ratesService.ErrInvalidAvgRange,
			ratesService.ErrNotEnoughValues,
			ratesService.ErrEmptyValues:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &ratesv1.GetRatesResponse{
		Ask:           rate.Ask,
		Bid:           rate.Bid,
		TimestampUnix: rate.ReceivedAt.Unix(),
	}, nil
}

func mapCalculationParams(params *ratesv1.CalculationParams) models.CalculationParams {
	if params == nil {
		return models.CalculationParams{}
	}

	return models.CalculationParams{
		Method: mapCalcMethod(params.GetMethod()),
		N:      int(params.GetN()),
		M:      int(params.GetM()),
	}
}

func mapCalcMethod(method ratesv1.CalcMethod) types.CalcMethod {
	switch method {
	case ratesv1.CalcMethod_CALC_METHOD_TOP_N:
		return types.CalcMethodTopN
	case ratesv1.CalcMethod_CALC_METHOD_AVG_NM:
		return types.CalcMethodAvgNM
	default:
		return ""
	}
}
