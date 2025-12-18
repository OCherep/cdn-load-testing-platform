package geo

import "net/http"

func Apply(req *http.Request, region string) {
	req.Header.Set("X-Geo-Region", region)
	req.Header.Set("X-Forwarded-For", fakeIP(region))
}

func fakeIP(region string) string {
	switch region {
	case "eu":
		return "2.16.0.1"
	case "us":
		return "23.0.0.1"
	case "asia":
		return "43.0.0.1"
	default:
		return "1.1.1.1"
	}
}
