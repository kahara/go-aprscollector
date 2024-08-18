package main

import (
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

func runFakeAprsisServer(addrPort string, term <-chan bool, ack chan<- bool) {
	l, err := net.Listen("tcp", addrPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not listen on " + addrPort)
	}
	defer l.Close()

	for {
		// Single connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not accept connection")
		}

		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-term:
				ack <- true
				return
			case <-ticker.C:
				conn.Write([]byte("Foo!\n"))
			}
		}
	}
}
