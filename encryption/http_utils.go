package main

import (
	"net/http"
	"time"
	"log"
	"net/http/httputil"
	"fmt"
)

// Intercepts existing round-tripper using logging implementation and uses
// default logger for request and response logging.
func NewDefaultLoggingTransport(rt http.RoundTripper) http.RoundTripper {
	return NewLoggingTransport(rt, DefaultLogger{})
}

// Intercepts existing round-tripper using logging implementation and uses
// provided logger for request and response logging.
func NewLoggingTransport(rt http.RoundTripper, logger HttpLogger) http.RoundTripper {
	return &loggingRoundTripper{
		roundTripper: rt,
		logger: logger,
	}
}

type (
	loggingRoundTripper struct {
		roundTripper http.RoundTripper
		logger       HttpLogger
	}

	// Interface defines method for logging requests and responses
	HttpLogger interface {
		LogRequest(*http.Request)
		LogResponse(*http.Request, *http.Response, time.Duration)
	}

	// Default implementation of logging request and response
	DefaultLogger struct {
	}
)

// Intercepts another RoundTripper implementation by request, response logging and duration calculation
func (c *loggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	c.logger.LogRequest(request)
	startTime := time.Now()
	response, err := c.roundTripper.RoundTrip(request)
	if err != nil {
		log.Fatalf("%s\n\n", err)
	} else {
		duration := time.Since(startTime)
		c.logger.LogResponse(request, response, duration)
	}
	return response, err
}

func (logger DefaultLogger) LogRequest(req *http.Request) {
	printDebug(httputil.DumpRequestOut(req, true))
}

func (logger DefaultLogger) LogResponse(req *http.Request, resp *http.Response, duration time.Duration) {
	printDebug(httputil.DumpResponse(resp, true))
}

func printDebug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}