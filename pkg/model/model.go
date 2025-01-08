package model

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"

	"github.com/pgillich/micro-server/pkg/configs"
)

type contextKey string

const (
	CtxKeyCmd              = contextKey("command")
	CtxKeyHttpServerRunner = contextKey("HttpServerController")
	CtxKeyBuildInfo        = contextKey("BuildInfo")
	CtxKeyServerConfig     = contextKey("ServerConfig")
	CtxKeyTestConfig       = contextKey("TestConfig")
)

// HttpServicer is the interface for HTTP services
type HttpServicer interface {
	// Name returns the name of the service
	Name() string
	// Prepare prepares the service for running, for example by registering HTTP routes, allocating resources, checking network dependencies, etc.
	Prepare(ctx context.Context, serverConfig configs.ServerConfiger, testConfig configs.TestConfiger, httpRouter chi.Router, tr trace.Tracer) error
	// Run runs the service. Called after all services have been prepared.
	Start(ctx context.Context) error
	// Stop stops the service. Called when the application is shutting down. It should free resources, close connections, etc.
	Stop(ctx context.Context) error
}

type BuildInfo interface {
	Version() string
	BuildTime() string
	AppName() string
	ModulePath() string
}
