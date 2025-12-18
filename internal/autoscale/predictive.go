package autoscale

import "math"

type Metrics struct {
	AvgLatency float64
	ErrorRate  float64
	RPS        int
}

type Decision struct {
	ScaleRPS   int
	ScaleNodes int
	Reason     string
}

type Predictor struct {
	TargetLatency float64
	MaxErrorRate  float64
	MaxRPSPerNode int
}

func NewPredictor() *Predictor {
	return &Predictor{
		TargetLatency: 200,
		MaxErrorRate:  0.01,
		MaxRPSPerNode: 2000,
	}
}

func (p *Predictor) Decide(m Metrics, nodes int) Decision {
	decision := Decision{
		ScaleRPS:   m.RPS,
		ScaleNodes: nodes,
		Reason:     "stable",
	}

	/*
		LATENCY PRESSURE
	*/
	if m.AvgLatency > p.TargetLatency {
		factor := m.AvgLatency / p.TargetLatency
		newNodes := int(math.Ceil(float64(nodes) * factor))

		decision.ScaleNodes = newNodes
		decision.Reason = "latency_pressure"
	}

	/*
		ERROR PRESSURE
	*/
	if m.ErrorRate > p.MaxErrorRate {
		decision.ScaleNodes = nodes + 1
		decision.Reason = "error_pressure"
	}

	/*
		RPS CEILING
	*/
	maxRPS := decision.ScaleNodes * p.MaxRPSPerNode
	if m.RPS > maxRPS {
		decision.ScaleNodes++
		decision.Reason = "rps_ceiling"
	}

	return decision
}
