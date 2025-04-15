package http

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// NopProbeProvider provides a NopProbe.
var NopProbeProvider = fx.Decorate(func() Probe { return NopProbe{} })

// NopProbe is a probe that does nothing.
type NopProbe struct{}

// Middleware is a middleware that does nothing.
func (probe NopProbe) Middleware(ctx *fiber.Ctx) error {
	return ctx.Next()
}
