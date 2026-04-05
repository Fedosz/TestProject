package grinex

import "errors"

var (
	ErrEmptyResponse = errors.New("empty response from grinex")
	ErrEmptyAsks     = errors.New("empty asks from grinex")
	ErrEmptyBids     = errors.New("empty bids from grinex")
)
