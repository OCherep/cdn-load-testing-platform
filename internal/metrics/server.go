package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("[metrics] listening on :2112")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()
}
