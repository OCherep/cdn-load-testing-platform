package load

import "cdn-load-platform/internal/metrics"

type SLA struct {
	LatencyMs  float64
	ErrorRate  float64
	Stickiness float64
}

func CheckSLA(sla SLA) (bool, string) {
	snap := metrics.Snapshot()

	if snap.P95Latency > sla.LatencyMs {
		return true, "Latency SLA breached"
	}

	if snap.ErrorRate > sla.ErrorRate {
		return true, "Error rate SLA breached"
	}

	if snap.StickinessRatio < sla.Stickiness {
		return true, "Stickiness SLA breached"
	}

	return false, ""
}
