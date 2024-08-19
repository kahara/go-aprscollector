package main

import (
	"github.com/kahara/go-canner"
	"github.com/rs/zerolog/log"
	"net"
	"testing"
	"time"
)

var (
	fakeAprsisServerAddress net.Addr
)

func fakeAprsisServer(packets <-chan string, term <-chan bool, ack chan<- bool) {
	l, err := net.Listen("tcp", "")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not listen ¯\\_(?)_/¯")
	}
	fakeAprsisServerAddress = l.Addr()
	defer l.Close()

	for {
		// Single connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not accept connection")
		}

		for {
			select {
			case <-term:
				ack <- true // Notify back that were done
				return
			case packet := <-packets:
				packet = packet + "\n"
				conn.Write([]byte(packet))
			}
		}
	}
}

func Test(t *testing.T) {
	serverPackets := make(chan string)
	serverTerm := make(chan bool)
	serverAck := make(chan bool)
	go fakeAprsisServer(serverPackets, serverTerm, serverAck)

	for {
		time.Sleep(time.Millisecond)
		if fakeAprsisServerAddress != nil {
			break
		}
	}

	config := Config{
		AprsisServer: fakeAprsisServerAddress.String(),
		LoginLine:    "user N0CALL pass -1 vers go-aprscollector-test v0 filter u/APBM1D",
	}

	tests := []struct {
		name    string
		packets []string
	}{
		{
			name: "APBM1D",
			packets: []string{"2M0IBX>APBM1D,GB7GL,DMR*,qAR,GB7GL:@015857h2300.00N/11300.00E[000/000John",
				"DO1AV>APBM1D,DB0FS,DMR*,qAR,DB0FS:@191249h5334.23N/01003.29EU061/000Auch QRV auf 145.525",
				"IW9HRH-7>APBM1D,IW9HRH,DMR*,qAR,IW9HRH:@015857h3412.00N/10850.00E[000/000Eugenio",
				"M0VXT-7>APBM1D,DMR*,qAS,BM2002-10:=5345.63N/00241.48W[180/091Gareth, TETRA, GB7PU",
				"SP5VTV-9>APBM1D,DMR*,qAS,SR0DMR-10:@191251h5211.98N/02036.35Ev171/000Anytone D878UV"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			records := make(chan canner.Record)
			term := make(chan bool)
			ack := make(chan bool)
			t.Logf("Connecting to " + config.AprsisServer)
			go Collect(&config, records, term, ack)
			for _, packet := range test.packets {
				t.Logf("Requesting packet '%s'", packet)
				serverPackets <- packet
				time.Sleep(10 * time.Millisecond)

				record := <-records
				t.Logf("Received record '%#v'", record)
			}
			term <- true
			<-ack
		})
	}

	// Shut the "fake" server down cleanly
	serverTerm <- true
	<-serverAck
}
