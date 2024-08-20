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

	// Collect
	collected := make(chan canner.Record, 1000)
	collectTerm := make(chan bool)
	collectAck := make(chan bool)
	go Collect(config, collected, collectTerm, collectAck)

	// Process
	processed := make(chan canner.Record, 1000)
	processTerm := make(chan bool)
	processAck := make(chan bool)
	go Process(config, collected, processed, processTerm, processAck)

	// Store

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	log.Info().Any("signal", <-sig).Msg("Signal caught, exiting")
}
