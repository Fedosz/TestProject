package rates

// calcTopN return top N value
func (s *Service) calcTopN(values []float64, n int) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptyValues
	}

	if n <= 0 {
		return 0, ErrInvalidTopN
	}

	if n > len(values) {
		return 0, ErrNotEnoughValues
	}

	return values[n-1], nil
}
