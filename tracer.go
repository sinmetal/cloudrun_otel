package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	cth, err := GetCloudTraceHeader(r)
	if err != nil {
		return ctx
	}
	traceID, err := trace.TraceIDFromHex(cth.TraceID)
	if err != nil {
		aelog.Warningf(ctx, "invalid TraceID err=%s", err)
		return ctx
	}
	spanID, err := trace.SpanIDFromHex(cth.SpanID)
	if err != nil {
		aelog.Warningf(ctx, "invalid SpanID err=%s", err)
		return ctx
	}
	return trace.ContextWithRemoteSpanContext(ctx, trace.SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
	})
}

type CloudTraceHeader struct {
	TraceID   string
	SpanID    string
	TraceTrue int // -1=未指定, 0=Traceしない, 1=Traceする
}

// GetCloudTraceHeader
// https://cloud.google.com/trace/docs/setup?hl=en#force-trace
func GetCloudTraceHeader(r *http.Request) (*CloudTraceHeader, error) {
	cth := &CloudTraceHeader{
		TraceTrue: -1,
	}
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	fl := strings.Split(traceHeader, "/")
	if len(fl) < 2 {
		return nil, fmt.Errorf("invalid Header")
	}
	cth.TraceID = fl[0]
	sl := strings.Split(fl[1], ";")
	cth.SpanID = sl[0]
	if len(sl) < 2 {
		if sl[1] == "o=1" {
			cth.TraceTrue = 1
		} else if sl[1] == "o=0" {
			cth.TraceTrue = 0
		}
	}
	return cth, nil
}
