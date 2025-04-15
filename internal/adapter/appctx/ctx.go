package appctx

import (
	"context"

	"go.uber.org/fx"
)

// Module provides a context that is canceled when the application is stopped.
var Module = fx.Module(
	"appcontext",
	fx.Provide(func(lc fx.Lifecycle) context.Context {
		ctx, cancel := context.WithCancel(context.Background())

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				cancel()
				return nil
			},
		})

		return ctx
	}),
)
