package report

import "time"
import "fmt"

type SLAEvidence struct {
	TestID    string
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

	LatencyOK    bool
	ErrorRateOK  bool
	StickinessOK bool
}

func (e *SLAEvidence) Evaluate() {
	e.LatencyOK = e.P95LatencyMs <= e.LatencySLAms
	e.ErrorRateOK = e.ErrorRate <= e.ErrorRateSLA
	e.StickinessOK = e.StickinessRatio >= e.StickinessSLA
}

func (e *SLAEvidence) Summary() string {
	e.Evaluate()

	result := "SLA EVIDENCE REPORT\n\n"

	result += "Test ID: " + e.TestID + "\n"
	result += "Target: " + e.TargetURL + "\n"
	result += "Start: " + e.StartTime.String() + "\n"
	result += "End:   " + e.EndTime.String() + "\n\n"

	result += "Metrics:\n"
	result += "- Avg Latency: " + format(e.AvgLatencyMs) + " ms\n"
	result += "- P95 Latency: " + format(e.P95LatencyMs) + " ms\n"
	result += "- Error Rate:  " + format(e.ErrorRate*100) + " %\n"
	result += "- Stickiness:  " + format(e.StickinessRatio*100) + " %\n\n"

	result += "SLA:\n"
	result += "- Latency SLA:    " + boolMark(e.LatencyOK) + "\n"
	result += "- Error Rate SLA: " + boolMark(e.ErrorRateOK) + "\n"
	result += "- Stickiness SLA: " + boolMark(e.StickinessOK) + "\n"

	return result
}

func boolMark(ok bool) string {
	if ok {
		return "OK"
	}
	return "BREACH"
}

func format(v float64) string {
	return fmt.Sprintf("%.2f", v)
}
