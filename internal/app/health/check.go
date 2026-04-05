package health

import "context"

// Check checks service availability
func (s *Service) Check(ctx context.Context) error {
	return s.pinger.Ping(ctx)
}
