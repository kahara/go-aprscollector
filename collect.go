package main

import (
	"bufio"
	"fmt"
	"github.com/kahara/go-canner"
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

func Collect(config *Config, packets chan<- canner.Record, term <-chan bool, ack chan<- bool) {

	var (
		err      error
		hostport = config.AprsisServer
		backoff  = backoffDefault
		conn     net.Conn
		n        int
		scanner  *bufio.Scanner
	)

	for {
		select {
		case <-term:
			ack <- true
			return
		default:
		}

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
			select {
			case <-term:
				ack <- true
				return
			default:
			}

			if scanner.Scan() {
				line := scanner.Bytes()
				now := time.Now().UTC()
				log.Debug().Time("received", now).Str("line", string(line)).Msg("Line scanned")
				lineCopy := make([]byte, len(line))
				copy(lineCopy, line)
				packets <- canner.Record{
					Timestamp:   now,
					Description: "aprsis-raw",
					Payload:     lineCopy,
				}
			} else {
				log.Error().Err(scanner.Err()).Msgf("Error while scanning lines from %s, sleeping for %fs before starting over", hostport, backoffWhenConnected.Seconds())
				time.Sleep(backoffWhenConnected)
				break
			}
		}
	}
}
