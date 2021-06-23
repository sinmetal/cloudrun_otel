package main

import (
	"context"
	"net/http"

	"github.com/vvakame/sdlog/aelog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer = otel.GetTracerProvider().Tracer("example.com/trace")

func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName)
}

func EndSpan(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

func SetAttributesKV(ctx context.Context, kv map[string]interface{}) {
	span := trace.SpanFromContext(ctx)
	for k, v := range kv {
		switch v := v.(type) {
		case string:
			span.SetAttributes(label.String(k, v))
		case bool:
			span.SetAttributes(label.Bool(k, v))
		case int:
			span.SetAttributes(label.Int(k, v))
		case int64:
			span.SetAttributes(label.Int64(k, v))
		case float32:
			span.SetAttributes(label.Float32(k, v))
		case float64:
			span.SetAttributes(label.Float64(k, v))
		default:
			span.SetAttributes(label.Any(k, v))
		}
	}
}

func SpanContextFromHttpRequest(ctx context.Context, r *http.Request) context.Context {
	traceIDHex := r.Header.Get("X-Cloud-Trace-Context")
	if traceIDHex == "" {
		return ctx
	}
	aelog.Infof(ctx, "TraceID:%s", traceIDHex)
	traceID, err := trace.TraceIDFromHex(traceIDHex)
	if err != nil {
		aelog.Warningf(ctx, "warning: failed err=%s", err)
		return ctx
	}
	return trace.ContextWithRemoteSpanContext(ctx, trace.SpanContext{
		TraceID: traceID,
	})
}
