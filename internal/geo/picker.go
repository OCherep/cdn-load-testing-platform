package geo

import "math/rand"

func Pick(distribution map[string]int) string {
	total := 0
	for _, w := range distribution {
		total += w
	}

	r := rand.Intn(total)
	acc := 0

	for region, w := range distribution {
		acc += w
		if r < acc {
			return region
		}
	}
	return "unknown"
}
