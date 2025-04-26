package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Module provides tracing functionality
func Module() fx.Option {
	return fx.Options(

		fx.Provide(
			func(cfg *Config) (*otlptrace.Exporter, error) {
				if !cfg.Enable {
					return nil, nil
				}

				opts := []otlptracegrpc.Option{
					otlptracegrpc.WithEndpoint(cfg.Addr),
					otlptracegrpc.WithTimeout(cfg.Timeout),
				}

				if !cfg.Secure {
					opts = append(opts, otlptracegrpc.WithInsecure())
				}

				return otlptracegrpc.NewUnstarted(opts...), nil
			},

			func(cfg *Config, exporter *otlptrace.Exporter) (*sdktrace.TracerProvider, error) {
				if !cfg.Enable {
					return nil, nil
				}

				opts := tracerProviderOpts(cfg, exporter)
				provider := sdktrace.NewTracerProvider(opts...)

				otel.SetTracerProvider(provider)

				otel.SetTextMapPropagator(
					propagation.NewCompositeTextMapPropagator(
						propagation.TraceContext{},
						propagation.Baggage{},
					),
				)

				return provider, nil
			},

			func(cfg *Config, provider *sdktrace.TracerProvider) (trace.Tracer, error) {
				if !cfg.Enable {
					return otel.Tracer("wasabi-disabled"), nil
				}

				return provider.Tracer("wasabi", tracerOpts(&cfg.Opts)...), nil
			},
		),

		fx.Invoke(
			func(
				lifecycle fx.Lifecycle,
				cfg *Config,
				exporter *otlptrace.Exporter,
				provider *sdktrace.TracerProvider,
			) {
				if !cfg.Enable {
					return
				}

				lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return exporter.Start(ctx)
					},
					OnStop: func(ctx context.Context) error {
						if err := exporter.Shutdown(ctx); err != nil {
							return err
						}
						return provider.Shutdown(ctx)
					},
				})
			},
		),
	)
}

func tracerProviderOpts(config *Config, exporter *otlptrace.Exporter) []sdktrace.TracerProviderOption {
	tpOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SampleRate)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.Opts.ServiceName),
		)),
	}

	return tpOpts
}
