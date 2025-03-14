package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/pgillich/micro-server/pkg/configs"
	"github.com/pgillich/micro-server/pkg/logger"
	"github.com/pgillich/micro-server/pkg/middleware"
	mw_server "github.com/pgillich/micro-server/pkg/middleware/server"
	"github.com/pgillich/micro-server/pkg/model"
	"github.com/pgillich/micro-server/pkg/tracing"
)

type ConfigSetter interface {
	SetListenAddr(string)
	SetInstance(string)
	SetCommand(string)
	SetJaegerURL(string)
	GetOptions() []string
}

var (
	ErrInvalidServerRunner = errors.New("invalid server runner")

	serviceFactories = map[string]func() model.HttpServicer{}
)

func RegisterHttpService(factory func() model.HttpServicer) {
	serviceFactories[factory().Name()] = factory
}

func RunHttpServer(h http.Handler, shutdown <-chan struct{}, addr string, log *slog.Logger) {
	server := &http.Server{ // nolint:gosec // not secure
		Handler: h,
		Addr:    addr,
	}

	go func() {
		<-shutdown
		if err := server.Shutdown(context.Background()); !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server shutdown error", logger.KeyError, err)
		}
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("Server exit error", logger.KeyError, err)
	} else {
		log.Info("Server exit")
	}
}

func RunServices(ctx context.Context, buildinfo model.BuildInfo, serviceNames []string,
	serverConfig configs.ServerConfiger, testConfig configs.TestConfiger,
) error {
	_, log := logger.FromContext(ctx)
	log.Info("SERVICES_TO_RUN", "services", strings.Join(serviceNames, ","))
	mux := chi.NewRouter()
	services := map[string]model.HttpServicer{}
	shutdown := make(chan struct{})

	for _, serviceName := range serviceNames {
		serviceFactory, exists := serviceFactories[serviceName]
		if !exists {
			return errors.New("service not found: " + serviceName)
		}
		services[serviceName] = serviceFactory()
	}

	for _, service := range services {
		log.Debug("SERVICE_PREPARE", logger.KeyService, service.Name())
		deferFn, err := PrepareService(ctx, service, buildinfo, serverConfig, testConfig, mux)
		defer deferFn()
		if err != nil {
			return err
		}
	}
	for _, service := range services {
		log.Info("SERVICE_START", logger.KeyService, service.Name())
		if err := service.Start(ctx); err != nil {
			return err
		}
	}

	httpServerRunner := RunHttpServer
	if testConfig.GetHttpServerRunner() != nil {
		httpServerRunner = testConfig.GetHttpServerRunner()
	}

	httpServerRunner(mux, shutdown, serverConfig.GetListenAddr(), log)
	log.Info("SERVER_STARTED")

	return nil
}

func PrepareService(ctx context.Context, service model.HttpServicer, buildinfo model.BuildInfo,
	serverConfig configs.ServerConfiger, testConfig configs.TestConfiger, mux *chi.Mux,
) (func(), error) {
	_, log := logger.FromContext(ctx)
	hostname, _ := os.Hostname() //nolint:errcheck // not important
	deferFn := func() {}

	traceOptions := []otlptracehttp.Option{}
	if serverConfig.GetTracerUrl() != "" {
		traceOptions = append(traceOptions, otlptracehttp.WithEndpointURL(serverConfig.GetTracerUrl()))
		traceOptions = append(traceOptions, otlptracehttp.WithCompression(otlptracehttp.NoCompression))
		traceOptions = append(traceOptions, otlptracehttp.WithInsecure())
	}
	traceExporter, err := tracing.OtlpProvider(ctx, traceOptions...)
	if err != nil {
		return deferFn, logger.Wrap(errors.New("unable to get OtlpProvider"), err)
	}
	tp := tracing.InitTracer(traceExporter, sdktrace.AlwaysSample(),
		buildinfo, service.Name(), hostname, "", log,
	)
	deferFn = func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("Error shutting down tracer provider", logger.KeyService, service.Name(), logger.KeyError, err)
		}
	}
	tr := tp.Tracer(
		buildinfo.ModulePath(),
		trace.WithInstrumentationVersion(tracing.SemVersion()),
	)

	r := mux.With(
		chi_middleware.Recoverer,
		mw_server.ChiLoggerBaseMiddleware(log.With(logger.KeyService, service.Name())),
		mw_server.ChiTracerMiddleware(tr, hostname),
		mw_server.ChiLoggerMiddleware(slog.LevelInfo, slog.LevelInfo),
		mw_server.ChiMetricMiddleware(middleware.GetMeter(buildinfo, log),
			"http_in", "HTTP in response", map[string]string{
				logger.KeyService: service.Name(),
			}, log,
		),
		chi_middleware.Recoverer,
	)

	if err := service.Prepare(ctx, serverConfig, testConfig, r, tr); err != nil {
		return deferFn, err
	}

	return deferFn, nil
}
