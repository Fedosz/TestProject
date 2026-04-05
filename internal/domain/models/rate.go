package models

import "time"

type Rate struct {
	Ask        float64
	Bid        float64
	ReceivedAt time.Time
}
