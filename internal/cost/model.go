package cost

type Estimate struct {
	HourlyUSD float64
}

func EstimateAgents(count int, pricePerHour float64) Estimate {
	return Estimate{
		HourlyUSD: float64(count) * pricePerHour,
	}
}
