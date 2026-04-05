package health

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"

	"rates_project/internal/app/health/mocks"
)

func TestService_Check(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(pinger *mocks.PingerMock)
		wantErr error
	}{
		{
			name: "success",
			prepare: func(pinger *mocks.PingerMock) {
				pinger.PingMock.Return(nil)
			},
		},
		{
			name: "ping error",
			prepare: func(pinger *mocks.PingerMock) {
				pinger.PingMock.Return(errors.New("ping error"))
			},
			wantErr: errors.New("ping error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			pinger := mocks.NewPingerMock(mc)

			if tt.prepare != nil {
				tt.prepare(pinger)
			}

			service := NewService(pinger)

			err := service.Check(context.Background())
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
