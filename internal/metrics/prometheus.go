package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

func RecordEdge(edge, host string, latency int64) {
	EdgeLatency.WithLabelValues(edge, host).Set(float64(latency))
}

func RecordStickiness(client string, ratio float64) {
	StickinessRatio.WithLabelValues(client).Set(ratio)
}
