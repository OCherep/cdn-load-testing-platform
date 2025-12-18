package metrics

import "github.com/prometheus/client_golang/prometheus"

var EdgeLatency = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "cdn_edge_latency_ms",
		Help: "Latency per CDN edge",
	},
	[]string{"edge", "region"},
)

func RegisterEdgeGeo() {
	prometheus.MustRegister(EdgeLatency)
}
