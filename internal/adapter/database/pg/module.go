package pg

import (
	"go.uber.org/fx"
)

var (
	Module = fx.Module("postgres", Providers, Invokers)

	Providers = fx.Options(
		fx.Provide(NewGorm),
		fx.Provide(NewSQLDB),
		fx.Provide(NewHealthChecker),
	)

	Invokers = fx.Options(
		fx.Invoke(EnableTracing),
		fx.Invoke(HookConnection),
		fx.Invoke(configureConnPool),
	)
)
