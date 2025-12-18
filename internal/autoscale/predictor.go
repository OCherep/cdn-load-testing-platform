package autoscale

import "cdn-load-platform/internal/metrics"

type Predictor struct {
	window []metrics.Snapshot
	size   int
}

func NewPredictor(size int) *Predictor {
	return &Predictor{size: size}
}

func (p *Predictor) Add(s metrics.Snapshot) {
	p.window = append(p.window, s)
	if len(p.window) > p.size {
		p.window = p.window[1:]
	}
}

func (p *Predictor) TrendRPS() int {
	if len(p.window) < 2 {
		return p.window[len(p.window)-1].RPS
	}
	delta := p.window[len(p.window)-1].RPS - p.window[0].RPS
	return p.window[len(p.window)-1].RPS + delta/2
}
