package main

import (
	"github.com/ebarkie/aprs"
	"github.com/kahara/go-canner"
	"github.com/rs/zerolog/log"
)

func Process(config *Config, packets <-chan canner.Record) {
	var frame aprs.Frame

	select {
	case packet := <-packets:
		log.Debug().Any("packet", string(packet.Payload)).Msg("Processing packet")

		if err := frame.FromBytes(packet.Payload); err != nil {
			log.Error().Err(err).Msg("Failed to process packet")
		} else {
			log.Debug().Any("frame", frame).Msg("")
		}
	}
}
