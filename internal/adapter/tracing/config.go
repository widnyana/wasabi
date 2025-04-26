package tracing

import (
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationVersion = "0.0.1"
	sdkVersion             = "1.30.0"
)

type (
	Config struct {
		Enable  bool          `mapstructure:"enable"`
		Name    string        `mapstructure:"name"`
		Addr    string        `mapstructure:"addr"`
		Secure  bool          `mapstructure:"secure"`
		Timeout time.Duration `mapstructure:"timeout"`
		Opts    TracerOptions `mapstructure:"opts"`
		// SampleRate defines the fraction of traces to be sampled.
		// Use a float between 0.0 (0%) and 1.0 (100%) to control trace volume.
		//
		// Recommended SampleRate values per environment:
		//
		//   Environment   | SampleRate     | Notes
		//   ------------- | -------------- | --------------------------------------------------------------
		//   Local dev     | 1.0            | Full trace visibility is helpful during development.
		//   Staging       | 0.5 - 1.0      | High visibility with moderate cost control.
		//   Production    | 0.01 - 0.2     | Balance cost and observability; 0.1 (10%) is common.
		//   High traffic  | 0.001 - 0.01   | Reduce telemetry backend load while retaining useful signals.
		SampleRate float64 `mapstructure:"sample_rate" default:"0.01"`
	}

	// TracerOptions contains configuration options for initializing the OpenTelemetry tracer.
	// These fields help enrich trace data with contextual information about the application and its environment.
	TracerOptions struct {
		// DeploymentEnv is the environment in which the application is running.
		// Common values include "local", "staging", or "production".
		// Example: "production"
		DeploymentEnv string `mapstructure:"deployment_env"`

		// DeploymentName identifies the specific deployment or cluster name.
		// This can be helpful when you have multiple deployments in different regions.
		// Example: "cluster-europe-west1"
		DeploymentName string `mapstructure:"deployment_name"`

		// InstrumentationName is the name of the module or library generating telemetry data.
		// This helps in identifying which instrumentation produced a span.
		// Example: "otel-http-server"
		InstrumentationName string `mapstructure:"instrumentation_name"`

		// InstrumentationVer is the semantic version of the instrumentation library.
		// Useful for tracing issues introduced by changes in observability instrumentation.
		// Example: "v1.5.0"
		InstrumentationVer string `mapstructure:"instrumentation_ver"`

		// ServiceName is the name of your application or service.
		// This is a required field for most trace backends and should be consistent.
		// Example: "auth-service"
		ServiceName string `mapstructure:"service_name"`

		// ServiceVersion indicates the current version of your service.
		// Important for tracking deployments and debugging version-specific issues.
		// Example: "1.2.3"
		ServiceVersion string `mapstructure:"service_version"`
	}
)

func tracerOpts(opts *TracerOptions) []trace.TracerOption {
	tOpts := []trace.TracerOption{
		trace.WithSchemaURL(semconv.SchemaURL),
		trace.WithInstrumentationVersion(instrumentationVersion),
		trace.WithInstrumentationAttributes(
			semconv.ServiceName(opts.ServiceName),
			semconv.ServiceVersionKey.String(opts.ServiceVersion),
			semconv.DeploymentNameKey.String(opts.DeploymentName),
			semconv.DeploymentEnvironmentNameKey.String(opts.DeploymentEnv),
			semconv.TelemetrySDKNameKey.String("opentelemetry"),
			semconv.TelemetrySDKVersionKey.String(sdkVersion),
		),
	}

	return tOpts
}
