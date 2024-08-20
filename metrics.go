package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	Namespace = "aprscollector"
	Subsystem = "aprsis"
)

var (
	packets_received *prometheus.CounterVec
	packets_skipped  *prometheus.CounterVec
)

func SetupMetrics() {
	packets_received = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "packets_received_total",
	}, []string{"tocall"})

	packets_skipped = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "packets_skipped_total",
	}, []string{"tocall"})
}

func Metrics(addrPort string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addrPort, nil); err != nil {
		log.Fatal().Err(err).Str("addrport", addrPort).Msg("Could not expose Prometheus metrics")
	}
}
