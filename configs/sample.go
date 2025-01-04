package configs

import (
	"github.com/pgillich/micro-server/pkg/configs"
)

//mw_client_model "github.com/pgillich/micro-server/pkg/middleware/client/model"

type ServerConfig struct {
	ListenAddr string
	TracerUrl  string
	Sample     SampleConfig
}

func (c *ServerConfig) GetListenAddr() string {
	return c.ListenAddr
}

func (c *ServerConfig) GetTracerUrl() string {
	return c.TracerUrl
}

type SampleConfig struct {
}

type TestConfig struct {
	/*
		CaptureTransportMode mw_client_model.CaptureTransportMode
		CaptureDir           string
		CaptureMatchers      []mw_client_model.CaptureMatcher
	*/
	HttpServerRunner configs.HttpServerRunner
}

func (c *TestConfig) GetHttpServerRunner() configs.HttpServerRunner {
	return c.HttpServerRunner
}

func (c *TestConfig) SetHttpServerRunner(httpServerRunner configs.HttpServerRunner) {
	c.HttpServerRunner = httpServerRunner
}
