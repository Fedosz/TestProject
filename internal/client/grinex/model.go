package grinex

type DepthItem struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
}

type DepthResponse struct {
	Timestamp int64       `json:"timestamp"`
	Asks      []DepthItem `json:"asks"`
	Bids      []DepthItem `json:"bids"`
}
