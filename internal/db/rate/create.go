package rate

import (
	"context"

	"rates_project/internal/domain/models"
	"rates_project/internal/telemetry"
)

// Create creates rates info row
func (r *Repo) Create(ctx context.Context, rate *models.Rate) error {
	ctx, span := telemetry.Tracer("rate_repo").Start(ctx, "rate_repo.Create")
	defer span.End()

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
