package geo

import "os"

func DetectRegion() Region {
	if r := os.Getenv("AGENT_REGION"); r != "" {
		return Region(r)
	}
	return EU
}
