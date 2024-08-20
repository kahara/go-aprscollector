package main

import (
	"github.com/kahara/go-canner"
	"github.com/rs/zerolog/log"
	"strings"
)

func isSamecall(packet []byte) bool {
	addressing := strings.Split(strings.ToUpper(string(packet)), ":")
	if len(addressing) < 2 {
		return false
	}
	addressing = strings.Split(addressing[0], ",")
	if len(addressing) < 2 {
		return false
	}
	callsign := strings.Split(addressing[0], ">")[0]
	callsign = strings.Split(callsign, "-")[0]

	for _, part := range addressing[1:] {
		if strings.Contains(part, callsign) {
			return true
		}
	}

	return false
}

func Process(config *Config, collected <-chan canner.Record, processed chan<- canner.Record, term <-chan bool, ack chan<- bool) {

	for {
		select {
		case <-term:
			log.Info().Msg("Terminating processing")
			ack <- true
			return
		case record := <-collected:
			if config.SkipSamecall && isSamecall(record.Payload) {
				log.Debug().Str("packet", string(record.Payload)).Msg("Skipping same call")
				continue
			}
			processed <- record
		}
	}
}
