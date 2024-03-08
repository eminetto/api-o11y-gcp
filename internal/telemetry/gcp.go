package telemetry

import (
	"context"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"go.opentelemetry.io/otel/trace"
)

// GCPOTel struct
type GCPOTel struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

// NewGCP instantiate a new GCPOtel struct
func NewGCP(ctx context.Context, serviceName string) (*GCPOTel, error) {
	var tp *sdktrace.TracerProvider
	var err error
	tp, err = createOtelTraceProvider(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := tp.Tracer(serviceName)

	return &GCPOTel{
		provider: tp,
		tracer:   tracer,
	}, nil
}

// Start a trace
func (ot *GCPOTel) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	if len(opts) == 0 {
		return ot.tracer.Start(ctx, name)
	}
	return ot.tracer.Start(ctx, name, opts[0])
}

// Shutdown finalize a trace
func (ot *GCPOTel) Shutdown(ctx context.Context) {
	ot.provider.Shutdown(ctx)
}

func createOCGPtelTraceProvider(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return nil, err
	}

	// Identify your application using resource detection
	res, err := resource.New(ctx,
		// Use the GCP resource detector to detect information about the GCP platform
		resource.WithDetectors(gcp.NewDetector()),
		// Keep the default detectors
		resource.WithTelemetrySDK(),
		// Add your own custom attributes to identify your application
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}
