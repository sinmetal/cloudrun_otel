package main

import (
	"io"
	"log"
	"net/http"
	"os"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func initTracer() func() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Create Google Cloud Trace exporter to be able to retrieve
	// the collected spans.
	_, flush, err := cloudtrace.InstallNewPipeline(
		[]cloudtrace.Option{cloudtrace.WithProjectID(projectID)},
		// For this example code we use sdktrace.AlwaysSample sampler to sample all traces.
		// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return flush
}

func main() {
	flush := initTracer()
	defer flush()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(label.String("server", "handling this..."))

		_, _ = io.WriteString(w, "Hello, world!\n")
	}
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")
	http.Handle("/hello", otelHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
