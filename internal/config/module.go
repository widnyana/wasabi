package config

import (
	"go.uber.org/fx"

	"github.com/widnyana/wasabi/internal/adapter/database/pg"
	"github.com/widnyana/wasabi/internal/adapter/http"
	"github.com/widnyana/wasabi/internal/adapter/logger"
	"github.com/widnyana/wasabi/internal/adapter/metrics"
	"github.com/widnyana/wasabi/internal/adapter/redis"
	"github.com/widnyana/wasabi/internal/adapter/tracing"
	"github.com/widnyana/wasabi/internal/constant"
)

var (
	// Module exports dependency to container
	Module = fx.Module("config", Providers, Invokers)

	// Invokers exports invokers to container
	Invokers = fx.Options(
		fx.Invoke(PrintBanner),
		fx.Invoke(LoadConfig),
	)

	// Providers exports providers to container
	// each provider will be called by fx.Options
	Providers = fx.Options(
		fx.Provide(NewAppConfig),
		fx.Provide(func(config *AppConfig) *http.Config { return &config.HTTP }),
		fx.Provide(func(config *AppConfig) *redis.Config { return &config.Redis }),
		fx.Provide(func(config *AppConfig) *pg.Config { return &config.Postgres }),
		fx.Provide(func(config *AppConfig) *metrics.Config { return &config.Metrics }),
		fx.Provide(func(config *AppConfig) *tracing.Config {
			cfg := &config.Tracing
			cfg.Name = constant.AppName
			cfg.Opts.InstrumentationName = constant.AppName
			cfg.Opts.ServiceVersion = constant.AppVersion

			return cfg
		}),
		fx.Provide(func(config *AppConfig) *logger.Config { return &config.Log }),
	)
)
