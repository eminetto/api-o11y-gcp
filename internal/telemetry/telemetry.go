package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Telemetry interface
type Telemetry interface {
	Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	Shutdown(ctx context.Context)
}
