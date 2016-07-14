package main

import (
	"net/http"
	"time"
	"net"
)

func init() {
	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, connectionTimeout)
		},
	}

	// Initialize HTTP client, if application contains -debug option
	// HTTP client will be initialized with logging transport to print
	// HTTP requests and responses.
	if debug {
		loggingTransport := NewDefaultLoggingTransport(&transport)
		httpClient = &http.Client{Transport: loggingTransport}
	} else {
		httpClient = &http.Client{Transport: &transport}
	}
}

const (
	// Settings used for all HTTP invocations
	connectionTimeout = time.Duration(2 * time.Second)
)

var httpClient *http.Client