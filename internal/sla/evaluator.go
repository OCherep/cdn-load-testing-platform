package sla

import "cdn-load-platform/internal/report"

func Evaluate(e *report.SLAEvidence) {
	e.LatencyBreached = e.AvgLatencyMs > e.LatencySLAms
	e.ErrorRateBreached = e.ErrorRate > e.ErrorRateSLA
	e.StickinessBreached = e.StickinessRatio < e.StickinessSLA
}

func IsBreached(e report.SLAEvidence) bool {
	return e.LatencyBreached || e.ErrorRateBreached || e.StickinessBreached
}
