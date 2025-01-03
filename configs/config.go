package configs

import (
	"log/slog"
	"net/http"
	//mw_client_model "github.com/pgillich/micro-server/pkg/middleware/client/model"
)

const (
	TracerVersion = "0.1.0"
)

type ServerConfig struct {
	ListenAddr string
	TracerUrl  string
	Sample     SampleConfig
}

type SampleConfig struct {
}

type TestConfig struct {
	/*
		CaptureTransportMode mw_client_model.CaptureTransportMode
		CaptureDir           string
		CaptureMatchers      []mw_client_model.CaptureMatcher
	*/
	HttpServerRunner HttpServerRunner
}

type HttpServerRunner func(h http.Handler, shutdown <-chan struct{}, addr string, l *slog.Logger)
