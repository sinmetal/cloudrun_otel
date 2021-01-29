package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer = otel.GetTracerProvider().Tracer("example.com/trace")

func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName)
}
