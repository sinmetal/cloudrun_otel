package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/datastore"
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/google/uuid"
	octrace "go.opencensus.io/trace"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	bridge "go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var ds *datastore.Client

func initTracer(projectID string) func() {
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
	octrace.DefaultTracer = bridge.NewTracer(otel.GetTracerProvider().Tracer("opencensus-bridge"))
	return flush
}

func initClient(ctx context.Context, projectID string) error {
	var err error
	ds, err = datastore.NewClient(ctx, projectID)
	return err
}

func main() {
	ctx := context.Background()
	var err error

	var projectID string
	if !metadata.OnGCE() {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	} else {
		projectID, err = metadata.ProjectID()
		if err != nil {
			log.Fatal(err)
		}
	}

	flush := initTracer(projectID)
	defer flush()

	if err := initClient(ctx, projectID); err != nil {
		log.Fatal(err)
	}
	als, err := NewAccessLogStore(ctx, ds)
	if err != nil {
		log.Fatal(err)
	}

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(label.String("server", "handling this..."))

		_, err := als.Insert(ctx, &AccessLog{
			ID: uuid.New().String(),
		})
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
		}

		_, _ = io.WriteString(w, "Hello, world!\n")
	}
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")
	http.Handle("/hello", otelHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
