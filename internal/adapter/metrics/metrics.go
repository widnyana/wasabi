package metrics

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ReadHeaderTimeout is the maximum amount of time to allow reading request headers.
const ReadHeaderTimeout = 5

type (
	// Config is the configuration for the metrics server.
	Config struct {
		Addr string `mapstructure:"addr"`
	}

	// Server is an alias for http.Server.
	Server = http.Server
)

// Module is the fx module for the metrics server.
var Module = fx.Module(
	"metrics",
	fx.Provide(NewServer),
	fx.Invoke(HookMetricsHandler),
)

// NewServer creates a new metrics server.
func NewServer(c Config) *Server {
	http.Handle("/metrics", promhttp.Handler())
	return &http.Server{
		Addr:              c.Addr,
		ReadHeaderTimeout: ReadHeaderTimeout * time.Second,
	}
}

// HookMetricsHandler hooks the metrics server to the fx lifecycle.
func HookMetricsHandler(lifecycle fx.Lifecycle, server *Server, logger *otelzap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				err := server.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("failed to start metrics server", zap.Error(err))
				}
			}()
			logger.Info("metrics server started", zap.String("addr", server.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
