package rate

import (
	"context"

	"rates_project/internal/telemetry"
)

// Ping pings database availability
func (r *Repo) Ping(ctx context.Context) error {
	ctx, span := telemetry.Tracer("rate_repo").Start(ctx, "rate_repo.Ping")
	defer span.End()

	return r.db.PingContext(ctx)
}
