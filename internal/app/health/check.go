package health

import "context"

func (s *Service) Check(ctx context.Context) error {
	return s.pinger.Ping(ctx)
}
