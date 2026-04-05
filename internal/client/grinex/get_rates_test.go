package grinex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_GetRates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		response    string
		statusCode  int
		wantAsks    []float64
		wantBids    []float64
		wantErrText string
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response: `{
				"timestamp": 1775419112,
				"asks": [
					{"price": "80.82", "volume": "10", "amount": "808.2"},
					{"price": "80.83", "volume": "11", "amount": "889.13"}
				],
				"bids": [
					{"price": "80.71", "volume": "12", "amount": "968.52"},
					{"price": "80.70", "volume": "13", "amount": "1049.10"}
				]
			}`,
			wantAsks: []float64{80.82, 80.83},
			wantBids: []float64{80.71, 80.70},
		},
		{
			name:       "empty response",
			statusCode: http.StatusOK,
			response: `{
				"timestamp": 1775419112,
				"asks": [],
				"bids": []
			}`,
			wantErrText: ErrEmptyResponse.Error(),
		},
		{
			name:       "empty asks",
			statusCode: http.StatusOK,
			response: `{
				"timestamp": 1775419112,
				"asks": [],
				"bids": [
					{"price": "80.71", "volume": "12", "amount": "968.52"}
				]
			}`,
			wantErrText: ErrEmptyAsks.Error(),
		},
		{
			name:       "empty bids",
			statusCode: http.StatusOK,
			response: `{
				"timestamp": 1775419112,
				"asks": [
					{"price": "80.82", "volume": "10", "amount": "808.2"}
				],
				"bids": []
			}`,
			wantErrText: ErrEmptyBids.Error(),
		},
		{
			name:       "invalid ask price",
			statusCode: http.StatusOK,
			response: `{
				"timestamp": 1775419112,
				"asks": [
					{"price": "invalid", "volume": "10", "amount": "808.2"}
				],
				"bids": [
					{"price": "80.71", "volume": "12", "amount": "968.52"}
				]
			}`,
			wantErrText: "invalid syntax",
		},
		{
			name:       "invalid bid price",
			statusCode: http.StatusOK,
			response: `{
				"timestamp": 1775419112,
				"asks": [
					{"price": "80.82", "volume": "10", "amount": "808.2"}
				],
				"bids": [
					{"price": "invalid", "volume": "12", "amount": "968.52"}
				]
			}`,
			wantErrText: "invalid syntax",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/spot/depth" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := NewClient(server.URL, "/api/v1/spot/depth", 5*time.Second)

			asks, bids, err := client.GetRates(context.Background())
			if tt.wantErrText != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErrText, err.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFloatSlicesEqual(t, tt.wantAsks, asks)
			assertFloatSlicesEqual(t, tt.wantBids, bids)
		})
	}
}

func assertFloatSlicesEqual(t *testing.T, expected, actual []float64) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("expected len %d, got %d", len(expected), len(actual))
	}

	for i := range expected {
		assertFloatEquals(t, expected[i], actual[i])
	}
}

func assertFloatEquals(t *testing.T, expected, actual float64) {
	t.Helper()

	const delta = 0.000001

	diff := expected - actual
	if diff < 0 {
		diff = -diff
	}

	if diff > delta {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
