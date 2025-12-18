package load

import "math/rand"

func CanaryURL(blue, green string, percent int) string {
	if rand.Intn(100) < percent {
		return green
	}
	return blue
}
