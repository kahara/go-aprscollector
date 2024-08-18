package main

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

const (
	backoffDefault       = 1.0
	backoffMultiplier    = 1.35
	backoffLimit         = 20
	backoffWhenConnected = time.Duration(5 * time.Second) // Sleep fixed amount of time when errors happen after connecting
)

func Collect(config *Config, packets chan []byte) {

	var (
		err      error
		hostport = config.AprsisServer
		backoff  = backoffDefault
		conn     net.Conn
		n        int
		scanner  *bufio.Scanner
	)

	for {
		if conn, err = net.Dial("tcp", hostport); err != nil {
			log.Error().Err(err).Msgf("Error while dialing %s, sleeping for %fs before starting over", hostport, backoff)
			time.Sleep(time.Duration(backoff) * time.Second)

			// Sleep longer after next error
			backoff = backoff * backoffMultiplier
			if backoff > backoffLimit {
				backoff = backoffLimit
			}
			continue
		}

		log.Printf("Connected to %s", hostport)
		backoff = backoffDefault

		// Log in
		if n, err = fmt.Fprint(conn, config.LoginLine); err != nil {
			log.Error().Err(err).Msgf("Error while logging in to %s, %d/%d bytes written; sleeping for %fs before starting over", hostport, n, len(config.LoginLine), backoffWhenConnected.Seconds())
			time.Sleep(backoffWhenConnected)
			break
		}

		// Consume APRS packets
		scanner = bufio.NewScanner(conn)
		for {
			if scanner.Scan() {
				line := scanner.Bytes()
				log.Debug().Str("line", string(line)).Msg("Line scanned")
				packets <- line
			} else {
				log.Error().Err(scanner.Err()).Msgf("Error while scanning lines from %s, sleeping for %fs before starting over", hostport, backoffWhenConnected.Seconds())
				time.Sleep(backoffWhenConnected)
				break
			}
		}
	}
}
