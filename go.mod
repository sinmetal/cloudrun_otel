module github.com/sinmetal/cloudrun_otel

go 1.15

require (
	cloud.google.com/go v0.75.0
	cloud.google.com/go/datastore v1.4.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.15.0
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/uuid v1.1.2
	github.com/vvakame/sdlog v0.0.0-20200409072131-7c0d359efddc
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.15.0
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/bridge/opencensus v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
)
