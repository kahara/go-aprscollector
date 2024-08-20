package main

import (
	"github.com/kahara/go-canner"
	"github.com/rs/zerolog/log"
)

func Store(config *Config, processed <-chan canner.Record, term <-chan bool, ack chan<- bool) {

	for {
		select {
		case <-term:
			log.Info().Msg("Terminating storage")
			ack <- true
			return
		case record := <-processed:
			log.Debug().Any("record", record).Msg("Storing record")
		}
	}
}
