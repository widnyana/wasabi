package http

import "go.uber.org/fx"

var (
	FiberModule = fx.Module("fiber", FiberProviders, FiberInvokes)

	FiberProviders = fx.Options(
		fx.Provide(NewFiber),
		fx.Provide(NewPromProbe),
		fx.Provide(func(probe *PromProbe) Probe { return probe }),
	)

	FiberInvokes = fx.Options(
		fx.Invoke(HookFiber),
	)
)
