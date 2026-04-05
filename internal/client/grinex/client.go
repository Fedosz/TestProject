package grinex

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
	path   string
}

func NewClient(baseURL, path string, timeout time.Duration) *Client {
	httpClient := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout)

	return &Client{
		client: httpClient,
		path:   path,
	}
}
