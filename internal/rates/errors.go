package rates

import "errors"

var (
	ErrInvalidCalcMethod = errors.New("invalid calc method")
	ErrInvalidTopN       = errors.New("invalid topN position")
	ErrInvalidAvgRange   = errors.New("invalid avgNM range")
	ErrNotEnoughValues   = errors.New("not enough values")
	ErrEmptyValues       = errors.New("empty values")
)
