package rates

import "testing"

func TestService_calcTopN(t *testing.T) {
	t.Parallel()

	service := &Service{}

	tests := []struct {
		name    string
		values  []float64
		n       int
		want    float64
		wantErr error
	}{
		{
			name:   "success first item",
			values: []float64{80.82, 80.83, 80.84},
			n:      1,
			want:   80.82,
		},
		{
			name:   "success third item",
			values: []float64{80.82, 80.83, 80.84},
			n:      3,
			want:   80.84,
		},
		{
			name:    "empty values",
			values:  []float64{},
			n:       1,
			wantErr: ErrEmptyValues,
		},
		{
			name:    "invalid n",
			values:  []float64{80.82},
			n:       0,
			wantErr: ErrInvalidTopN,
		},
		{
			name:    "not enough values",
			values:  []float64{80.82, 80.83},
			n:       3,
			wantErr: ErrNotEnoughValues,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := service.calcTopN(tt.values, tt.n)
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
