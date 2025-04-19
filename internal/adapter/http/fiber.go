package http

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	defaultIdleTimeoutSeconds = 5
	kilobyte                  = 1024
	fiberConcurrency          = 512 * kilobyte
)

// Probe is an interface for HTTP request probes.
type (
	Probe interface {
		Middleware(*fiber.Ctx) error
	}
)

// NewFiber creates a new Fiber app.
func NewFiber(logger *otelzap.Logger, probe Probe) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "wasabi",
		Concurrency:           fiberConcurrency,
		DisableStartupMessage: true,
		DisableKeepalive:      true,
		EnablePrintRoutes:     false,
		IdleTimeout:           time.Duration(defaultIdleTimeoutSeconds) * time.Second,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	})

	app.Use(otelfiber.Middleware(
		otelfiber.WithServerName("wasabi"),
		otelfiber.WithTracerProvider(otel.GetTracerProvider()),
	))
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Logger,
	}))

	app.Use(probe.Middleware)

	return app
}

// HookFiber hooks the Fiber app to the lifecycle.
func HookFiber(lifecycle fx.Lifecycle, app *fiber.App, config Config, logger *otelzap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				port := config.Port
				if port == 0 {
					port = 9999
				}

				addr := fmt.Sprintf(":%d", port)
				logger.Info("REST API Server running", zap.String("addr", addr))
				if err := app.Listen(addr); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})
}
