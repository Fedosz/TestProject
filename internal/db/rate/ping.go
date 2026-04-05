package rate

import "context"

func (r *Repo) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
