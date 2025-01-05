package sample

import (
	"context"
	"errors"

	"github.com/go-chi/chi/v5"

	"github.com/pgillich/micro-server/configs"
	srv_configs "github.com/pgillich/micro-server/pkg/configs"
	"github.com/pgillich/micro-server/pkg/model"
	"github.com/pgillich/micro-server/pkg/server"
)

const (
	ServiceName = "sample"
)

var (
	ErrUnableToPrepareService = errors.New("unable to prepare service")
)

type HttpService struct {
	serverConfig *configs.ServerConfig
	testConfig   *configs.TestConfig
	apiServer    *ApiServer
}

func newHttpService() model.HttpServicer {
	return &HttpService{}
}

func init() {
	server.RegisterHttpService(newHttpService)
}

func (s *HttpService) Name() string {
	return ServiceName
}

func (s *HttpService) Prepare(ctx context.Context, serverConfig srv_configs.ServerConfiger, testConfig srv_configs.TestConfiger,
	httpRouter chi.Router,
) error {
	var is bool
	s.serverConfig, is = serverConfig.(*configs.ServerConfig)
	if !is {
		return srv_configs.ErrFatalServerConfig
	}
	s.apiServer = &ApiServer{service: s}

	s.testConfig, is = testConfig.(*configs.TestConfig)
	if !is {
		return srv_configs.ErrFatalServerConfig
	}

	httpRouter.Get("/hello", s.apiServer.GetHello)

	return nil
}

func (s *HttpService) Start(ctx context.Context) error {
	return nil
}

func (s *HttpService) Stop(ctx context.Context) error {
	return nil
}
