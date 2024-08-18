package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	config := NewConfig()
	log.Info().Any("config", config).Msg("Configured")

	SetupMetrics()
	go Metrics(config.Metrics)

	packets := make(chan []byte, 1000)
	Collect(config, packets)
}
