package grinex

import (
	"context"
	"strconv"

	"rates_project/internal/telemetry"
)

// GetRates get current rates
func (c *Client) GetRates(ctx context.Context) ([]float64, []float64, error) {
	ctx, span := telemetry.Tracer("grinex_client").Start(ctx, "grinex.GetRates")
	defer span.End()

	var response DepthResponse

	_, err := c.client.R().
		SetContext(ctx).
		SetResult(&response).
		Get(c.path)
	if err != nil {
		return nil, nil, err
	}

	if len(response.Asks) == 0 && len(response.Bids) == 0 {
		return nil, nil, ErrEmptyResponse
	}

	if len(response.Asks) == 0 {
		return nil, nil, ErrEmptyAsks
	}

	if len(response.Bids) == 0 {
		return nil, nil, ErrEmptyBids
	}

	asks := make([]float64, 0, len(response.Asks))
	for _, item := range response.Asks {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return nil, nil, err
		}

		asks = append(asks, price)
	}

	bids := make([]float64, 0, len(response.Bids))
	for _, item := range response.Bids {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return nil, nil, err
		}

		bids = append(bids, price)
	}

	return asks, bids, nil
}
