package sla

type EdgeAffinityRule struct {
	MinRatio float64
}

func EdgeAffinityBreached(ratio float64, rule EdgeAffinityRule) bool {
	return ratio < rule.MinRatio
}
