package test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pgillich/micro-server/configs"
	"github.com/pgillich/micro-server/internal/buildinfo"
	_ "github.com/pgillich/micro-server/internal/sample"
	"github.com/pgillich/micro-server/pkg/logger"
	"github.com/pgillich/micro-server/pkg/testutil"
	// "github.com/pgillich/micro-server/internal/tracing"
)

type E2ETestSuite struct {
	suite.Suite
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) TestHello() {
	log := logger.GetLogger(buildinfo.BuildInfo.AppName(), slog.LevelDebug).With(logger.KeyTestCase, s.T().Name())
	//tracing.SetErrorHandlerLogger(log)
	serverConfig := &configs.ServerConfig{
		Sample: configs.SampleConfig{},
	}
	testConfig := &configs.TestConfig{}

	server := testutil.RunTestServerCmd(s.T(), "services", serverConfig, testConfig, []string{"sample"}, log)
	defer server.Cancel()

	testUrl, err := url.JoinPath(server.TestServer.URL, "/hello")
	s.NoError(err, "testUrl")
	clientCtx := logger.NewContext(context.Background(), log)

	req, err := http.NewRequestWithContext(clientCtx, http.MethodGet, testUrl, nil)
	s.NoError(err, "GetHello")
	client := http.Client{
		Transport: http.DefaultTransport,
	}
	resp, err := client.Do(req)
	s.NoError(err, "GetHello")
	s.NotNil(resp.Body, "GetHello")
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	s.NoError(err, "GetHello")
	bodyStr := string(body)
	s.Equal("World", bodyStr, "GetHello")

	s.T().Logf("Client Resp\n%s", bodyStr)

	//time.Sleep(1000 * time.Second)
}
