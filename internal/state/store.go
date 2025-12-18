package state

type ChaosConfig struct {
	Enabled    bool `json:"enabled"`
	LatencyMs  int  `json:"latency_ms"`
	ErrorRate  int  `json:"error_rate"`
	BurstPause bool `json:"burst_pause"`
}

type TestState struct {
	TestID string `json:"test_id"`
	Status string `json:"status"` // running | paused | stopped

	StartedAt int64 `json:"started_at"`
	ExpiresAt int64 `json:"expires_at"`

	MinRPS     int  `json:"min_rps"`
	MaxRPS     int  `json:"max_rps"`
	DesiredRPS *int `json:"desired_rps,omitempty"`

	ChaosConfig ChaosConfig `json:"chaos"`
}
