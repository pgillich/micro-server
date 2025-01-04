package testutil

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/pgillich/micro-server/internal/sample"
	"github.com/pgillich/micro-server/pkg/cmd"
	pkg_configs "github.com/pgillich/micro-server/pkg/configs"
	"github.com/pgillich/micro-server/pkg/logger"
	"github.com/pgillich/micro-server/pkg/utils"
	// "github.com/pgillich/micro-server/pkg/tracing"
)

type TestServer struct {
	TestServer *httptest.Server
	Addr       string
	Ctx        context.Context //nolint:containedctx // test
	Cancel     context.CancelFunc
}

func RunTestServerCmd(t *testing.T, serverName string,
	serverConfig pkg_configs.ServerConfiger, testConfig pkg_configs.TestConfiger,
	args []string, log *slog.Logger,
) *TestServer {
	ctx := logger.NewContext(context.Background(), log)

	server := &TestServer{
		TestServer: httptest.NewUnstartedServer(nil),
	}
	server.Addr = server.TestServer.Listener.Addr().String()

	started := make(chan struct{})
	runner := HttpTestserverRunner(server.TestServer, started)
	server.Ctx, server.Cancel = context.WithCancel(ctx)
	testConfig.SetHttpServerRunner(runner)

	workDir := t.TempDir()

	serverConfigFile, err := utils.SaveServerConfig(serverConfig, workDir, "micro_server.yaml")
	if err != nil {
		t.Fatalf("SaveServerConfig: %v", err)
	}

	command := append([]string{
		serverName,
		"--config", serverConfigFile,
		"--listenaddr", server.Addr,
	}, args...)

	go func() {
		cmd.Execute(server.Ctx, command, serverConfig, testConfig)
	}()
	<-started
	//time.Sleep(1 * time.Second)

	return server
}

func HttpTestserverRunner(server *httptest.Server, started chan struct{}) pkg_configs.HttpServerRunner {
	return func(h http.Handler, shutdown <-chan struct{}, addr string, log *slog.Logger) {
		server.Config.Handler = h
		log.Info("TestServer start")
		server.Start()
		close(started)
		log.Info("TestServer started")
		<-shutdown
		log.Info("TestServer shutdown")
		server.Close()
	}
}
