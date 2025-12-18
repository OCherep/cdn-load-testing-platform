package load

import (
	"sync"
	"time"
)

type Metrics struct {
	RPS       int
	P95       int
	ErrorRate float64
}

type AdaptiveEngine struct {
	mu         sync.Mutex
	currentRPS int
	maxRPS     int
	step       int
	latencySLA int
	errorSLA   float64
}

func NewAdaptiveEngine(start, max, step int) *AdaptiveEngine {
	return &AdaptiveEngine{
		currentRPS: start,
		maxRPS:     max,
		step:       step,
		latencySLA: 500,
		errorSLA:   1.0,
	}
}

func (e *AdaptiveEngine) Adjust(m Metrics) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch {
	case m.P95 > e.latencySLA || m.ErrorRate > e.errorSLA:
		e.currentRPS -= e.step
	case e.currentRPS < e.maxRPS:
		e.currentRPS += e.step
	}

	if e.currentRPS < 1 {
		e.currentRPS = 1
	}
	if e.currentRPS > e.maxRPS {
		e.currentRPS = e.maxRPS
	}

	return e.currentRPS
}

target := multiCDN.Pick(workerID)
req.URL = target.URL
metrics.RecordCDN(target.Name, latency)
