package health

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}

type Service struct {
	pinger Pinger
}

func NewService(pinger Pinger) *Service {
	return &Service{
		pinger: pinger,
	}
}
