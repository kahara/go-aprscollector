package main

import (
	"os"
	"strconv"
)

const (
	DefaultAprsisServer  = "rotate.aprs2.net:14580" // https://www.aprs-is.net/Connecting.aspx
	DefaultClientName    = "go-aprscollector"
	DefaultClientVersion = "v0"
	DefaultCallsign      = "N0CALL"
	DefaultPasscode      = "-1"
	DefaultFilter        = "u/APBM*" // https://github.com/aprsorg/aprs-deviceid/blob/ddfee32784215a2a4a5fe74c69812cfd87bd098c/tocalls.yaml#L341
	DefaultSkipSamecall  = true
	DefaultMetrics       = ":9108"
)

type Config struct {
	AprsisServer  string
	ClientName    string
	ClientVersion string
	Callsign      string
	Passcode      string
	Filter        string
	SkipSamecall  bool
	Metrics       string
}

func NewConfig() *Config {
	var config Config

	// APRS-IS server address:port
	aprsisServer := os.Getenv("APRSIS_SERVER")
	if aprsisServer == "" {
		config.AprsisServer = DefaultAprsisServer
	} else {
		config.AprsisServer = aprsisServer
	}

	// Client software name
	clientName := os.Getenv("CLIENT_NAME")
	if clientName == "" {
		config.ClientName = DefaultClientName
	} else {
		config.ClientName = clientName
	}

	// Client software version
	clientVersion := os.Getenv("CLIENT_VERSION")
	if clientVersion == "" {
		config.ClientVersion = DefaultClientVersion
	} else {
		config.ClientVersion = clientVersion
	}

	// Callsign
	callsign := os.Getenv("CALLSIGN")
	if callsign == "" {
		config.Callsign = DefaultCallsign
	} else {
		config.Callsign = callsign
	}

	// Passcode
	passcode := os.Getenv("PASSCODE")
	if passcode == "" {
		config.Passcode = DefaultPasscode
	} else {
		config.Passcode = passcode
	}

	// Filter
	filter := os.Getenv("FILTER")
	if filter == "" {
		config.Filter = DefaultFilter
	} else {
		config.Filter = filter
	}

	// Same call (i.e., destination contains source) skipping to e.g. ignore DMR hotspots' traffic
	skipSamecall := os.Getenv("SKIP_SAMECALL")
	if skipSamecall == "" {
		config.SkipSamecall = DefaultSkipSamecall
	} else {
		s, _ := strconv.ParseBool(skipSamecall)
		config.SkipSamecall = s
	}

	// Metrics address:port
	metrics := os.Getenv("METRICS")
	if metrics == "" {
		config.Metrics = DefaultMetrics
	} else {
		config.Metrics = metrics
	}

	return &config
}
