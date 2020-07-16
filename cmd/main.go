package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestNumber = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_number_total",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(requestNumber)
}

func main() {
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())

	http.DefaultClient.Timeout = 30 * time.Second

	http.HandleFunc("/data", func(rw http.ResponseWriter, r *http.Request) {
		requestNumber.Inc()
		rw.Write([]byte("OK"))
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	//Port 8080 to be redirected to Envoy proxy not reacheble by Prometheus
	log.Fatal(http.ListenAndServe(":8080", nil))
}
