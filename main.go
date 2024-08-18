package main

import (
	"github.com/kahara/go-canner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	config := NewConfig()
	log.Info().Any("config", config).Msg("Configured")

	SetupMetrics()
	go Metrics(config.Metrics)

	packets := make(chan canner.Record, 1000)
	go Collect(config, packets)
	go Process(config, packets)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	log.Info().Any("signal", <-sig).Msg("Signal caught, exiting")
}
