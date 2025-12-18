package state

type SLAConfig struct {
	LatencyMs  float64 `json:"latency_ms"`
	ErrorRate  float64 `json:"error_rate"`
	Stickiness float64 `json:"stickiness"`
}

type MetricsSnapshot struct {
	AvgLatency      float64
	P95Latency      float64
	ErrorRate       float64
	StickinessRatio float64
}

type TestState struct {
	TestID     string `json:"test_id"`
	Status     string `json:"status"`
	ProfileKey string `json:"profile_key"`

	Nodes    int `json:"nodes"`
	Sessions int `json:"sessions"`

	DesiredRPS int `json:"desired_rps"`

	StartedAt int64 `json:"started_at"`
	TTL       int64 `json:"ttl"`

	SLA SLAConfig `json:"sla"`

	SLABreached bool `json:"sla_breached"`
}
