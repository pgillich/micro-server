package configs

import (
	"errors"
	"log/slog"
	"net/http"
	//mw_client_model "github.com/pgillich/micro-server/pkg/middleware/client/model"
)

const (
	TracerVersion = "0.1.0"
)

var (
	ErrUnableToPrepareService = errors.New("unable to prepare service")
	ErrFatalServerConfig      = errors.New("fatal server config")
)

type ServerConfiger interface {
	GetListenAddr() string
	GetTracerUrl() string
}

type TestConfiger interface {
	GetHttpServerRunner() HttpServerRunner
	SetHttpServerRunner(HttpServerRunner)
}

type HttpServerRunner func(h http.Handler, shutdown <-chan struct{}, addr string, l *slog.Logger)