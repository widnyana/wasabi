package http

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,

	// TODO: register routing
	// fx.Invoke(func(app *fiber.App, r *root.Controller) { r.Register(app) }),
	// fx.Invoke(func(app *fiber.App, c *http.Controller) { c.Register(app) }),
)
