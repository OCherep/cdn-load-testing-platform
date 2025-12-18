package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Stickiness = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "cdn_stickiness_ratio",
		Help: "Edge stickiness ratio",
	},
)

func RecordStickiness(ratio float64) {
	Stickiness.Set(ratio)
}

var (
	RequestLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cdn_request_latency_ms",
			Help:    "Request latency",
			Buckets: prometheus.ExponentialBuckets(10, 2, 10),
		},
	)

	EdgeLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cdn_edge_latency_ms",
			Help: "Latency per edge",
		},
		[]string{"edge", "host"},
	)

	StickinessRatio = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cdn_edge_stickiness_ratio",
			Help: "Edge stickiness ratio per client",
		},
		[]string{"client"},
	)

	Latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cdn_latency_ms",
			Help:    "CDN request latency",
			Buckets: prometheus.ExponentialBuckets(10, 2, 10),
		},
		[]string{"cdn", "edge", "region"},
	)
)

func Start() {
	prometheus.MustRegister(RequestLatency)
	prometheus.MustRegister(EdgeLatency)
	prometheus.MustRegister(StickinessRatio)

	go http.ListenAndServe(":2112", promhttp.Handler())
}

func RecordLatency(ms int64) {
	RequestLatency.Observe(float64(ms))
}

func RecordLatency(cdn, edge, region string, value int64) {
	Latency.WithLabelValues(cdn, edge, region).
		Observe(float64(value))
}

func RecordStickiness(client string, ratio float64) {
	StickinessRatio.WithLabelValues(client).Set(ratio)
}

var Errors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cdn_errors_total",
		Help: "Total request errors",
	},
	[]string{"cdn", "region"},
)

func RecordError(cdn, region string) {
	Errors.WithLabelValues(cdn, region).Inc()
}
