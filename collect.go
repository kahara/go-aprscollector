package main

import (
	"github.com/rs/zerolog/log"
	"time"
)

func Collect(packets chan []byte) {
	for {
		time.Sleep(time.Second)
		log.Debug().Msg("Foo!")
	}
}
