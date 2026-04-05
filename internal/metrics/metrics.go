package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	GetRatesRequestsTotal prometheus.Counter
	GetRatesErrorsTotal   prometheus.Counter
	GetRatesDuration      prometheus.Histogram
	LastAsk               prometheus.Gauge
	LastBid               prometheus.Gauge
}

func MustNew() *Metrics {
	m := &Metrics{
		GetRatesRequestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "rates_get_rates_requests_total",
			Help: "Total number of GetRates requests",
		}),
		GetRatesErrorsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "rates_get_rates_errors_total",
			Help: "Total number of GetRates errors",
		}),
		GetRatesDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "rates_get_rates_duration_seconds",
			Help:    "GetRates request duration in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		LastAsk: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rates_last_ask",
			Help: "Last successful ask rate",
		}),
		LastBid: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rates_last_bid",
			Help: "Last successful bid rate",
		}),
	}

	prometheus.MustRegister(
		m.GetRatesRequestsTotal,
		m.GetRatesErrorsTotal,
		m.GetRatesDuration,
		m.LastAsk,
		m.LastBid,
	)

	return m
}
