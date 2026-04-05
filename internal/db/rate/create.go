package rate

import (
	"context"

	"rates_project/internal/domain/models"
)

func (r *Repo) Create(ctx context.Context, rate *models.Rate) error {
	query := `
		INSERT INTO rates (
			ask,
			bid,
			received_at
		) VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		rate.Ask,
		rate.Bid,
		rate.ReceivedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
