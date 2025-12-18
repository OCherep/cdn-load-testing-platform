package qoe

func Score(latencyMs float64, errorRate float64) float64 {
	score := 100.0

	if latencyMs > 300 {
		score -= (latencyMs - 300) * 0.1
	}
	score -= errorRate * 50

	if score < 0 {
		return 0
	}
	return score
}
