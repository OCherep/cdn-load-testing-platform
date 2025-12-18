package cost

func Estimate(nodes int, instance string, hours float64) float64 {
	price := EC2Hourly[instance]
	return float64(nodes) * price * hours
}
