package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const (
	Namespace = "aprscollector"
	Subsystem = "aprsis"
)

var (
	packets_received *prometheus.CounterVec
	packets_skipped  *prometheus.CounterVec
	packets_stored   *prometheus.CounterVec
)

func SetupMetrics() {
	packets_received = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "packets_received_total",
	}, []string{})

	packets_skipped = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "packets_skipped_total",
	}, []string{})
	packets_stored = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "packets_stored_total",
	}, []string{})
}

func Metrics(addrPort string) {

	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
	metricsServer := &http.Server{
		Addr:              addrPort,
		Handler:           metricsMux,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	log.Fatal().Err(metricsServer.ListenAndServe()).Send()
}
