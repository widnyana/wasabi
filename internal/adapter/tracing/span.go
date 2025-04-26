package tracing

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// https://pkg.go.dev/go.opentelemetry.io/otel/trace
type (
	Span   = trace.Span
	Tracer = trace.Tracer

	tracerKey struct{}
)

var global = noop.NewTracerProvider().Tracer("")

// FromContext returns the Tracer stored by a prior call to [WithTracer] or [SetDefaultTracer].
func FromContext(ctx context.Context) Tracer {
	if t, ok := ctx.Value(tracerKey{}).(Tracer); ok {
		return t
	}
	return global
}

// Start creates a Span and a Context containing it. It uses the Tracer returned by [FromContext].
func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	return FromContext(ctx).Start(ctx, name, opts...)
}

// Bool sets the k attribute of s to v.
func Bool(s Span, k string, v bool) { s.SetAttributes(attribute.Bool(k, v)) }

// Int sets the k attribute of s to v.
func Int(s Span, k string, v int) { s.SetAttributes(attribute.Int(k, v)) }

// String sets the k attribute of s to v.
func String(s Span, k, v string) { s.SetAttributes(attribute.String(k, v)) }
