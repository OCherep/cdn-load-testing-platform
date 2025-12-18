package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cdn_latency_ms",
			Help:    "CDN request latency",
			Buckets: prometheus.ExponentialBuckets(10, 2, 10),
		},
		[]string{"edge", "region"},
	)

	Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cdn_errors_total",
			Help: "Total request errors",
		},
		[]string{"region"},
	)
)

func RecordLatency(edge, region string, value int64) {
	Latency.WithLabelValues(edge, region).Observe(float64(value))
}

func RecordError(region string) {
	Errors.WithLabelValues(region).Inc()
}
