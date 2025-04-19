package config

import (
	"github.com/widnyana/wasabi/internal/adapter/database/pg"
	"github.com/widnyana/wasabi/internal/adapter/http"
	"github.com/widnyana/wasabi/internal/adapter/logger"
	"github.com/widnyana/wasabi/internal/adapter/metrics"
	"github.com/widnyana/wasabi/internal/adapter/redis"
	"github.com/widnyana/wasabi/internal/adapter/tracing"
	"go.uber.org/fx"
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
		fx.Provide(func(config *AppConfig) http.Config { return config.HTTP }),
		fx.Provide(func(config *AppConfig) redis.Config { return config.Redis }),
		fx.Provide(func(config *AppConfig) pg.Config { return config.Postgres }),
		fx.Provide(func(config *AppConfig) metrics.Config { return config.Metrics }),
		fx.Provide(func(config *AppConfig) tracing.Config { return config.Tracing }),
		fx.Provide(func(config *AppConfig) logger.Config { return config.Log }),
	)
)
