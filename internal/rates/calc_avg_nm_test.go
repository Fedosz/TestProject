package rates

import "testing"

func TestService_calcAvgNM(t *testing.T) {
	t.Parallel()

	service := &Service{}

	tests := []struct {
		name    string
		values  []float64
		n       int
		m       int
		want    float64
		wantErr error
	}{
		{
			name:   "success average first three",
			values: []float64{80.0, 81.0, 82.0},
			n:      1,
			m:      3,
			want:   81.0,
		},
		{
			name:   "success average middle range",
			values: []float64{80.0, 81.0, 82.0, 83.0},
			n:      2,
			m:      3,
			want:   81.5,
		},
		{
			name:    "empty values",
			values:  []float64{},
			n:       1,
			m:       2,
			wantErr: ErrEmptyValues,
		},
		{
			name:    "invalid range because n is zero",
			values:  []float64{80.0, 81.0},
			n:       0,
			m:       2,
			wantErr: ErrInvalidAvgRange,
		},
		{
			name:    "invalid range because n greater than m",
			values:  []float64{80.0, 81.0},
			n:       3,
			m:       2,
			wantErr: ErrInvalidAvgRange,
		},
		{
			name:    "not enough values",
			values:  []float64{80.0, 81.0},
			n:       1,
			m:       3,
			wantErr: ErrNotEnoughValues,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := service.calcAvgNM(tt.values, tt.n, tt.m)
			if err != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}

				return
			}

			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
