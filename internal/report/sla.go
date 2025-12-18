package report

import "time"

type SLAEvidence struct {
	TestID string

	TargetURL string
	StartTime time.Time
	EndTime   time.Time

	AvgLatencyMs    float64
	P95LatencyMs    float64
	ErrorRate       float64
	StickinessRatio float64

	LatencySLAms  float64
	ErrorRateSLA  float64
	StickinessSLA float64

	LatencyBreached    bool
	ErrorRateBreached  bool
	StickinessBreached bool
}
