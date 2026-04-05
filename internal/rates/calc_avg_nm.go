package rates

func (s *Service) calcAvgNM(values []float64, n, m int) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptyValues
	}

	if n <= 0 || m <= 0 || n > m {
		return 0, ErrInvalidAvgRange
	}

	if m > len(values) {
		return 0, ErrNotEnoughValues
	}

	var sum float64

	for i := n - 1; i <= m-1; i++ {
		sum += values[i]
	}

	count := m - n + 1

	return sum / float64(count), nil
}
