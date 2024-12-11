package otelinit

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

func InitTracing(conf Config, serviceName, packagePath string) (trace.Tracer, func(), error) {
	ctx := context.Background()

	exp, err := newTracingExporter(ctx, conf.ExporterType)
	if err != nil {
		return nil, nil, err
	}

	tp := newTraceProvider(serviceName, exp)

	otel.SetTracerProvider(tp)

	if conf.RegisterGlobal {
		Tracer = tp.Tracer(packagePath)
	}

	return Tracer, func() { _ = tp.Shutdown(ctx) }, nil
}

func newTracingExporter(_ context.Context, et ExporterType) (sdktrace.SpanExporter, error) {
	switch et {
	case Stdout:
		return stdouttrace.New()
	case Noop:
		return ExporterNoopTracing{}, nil
	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", et)
	}
}

func newTraceProvider(serviceName string, exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

type ExporterNoopTracing struct{}

func (e ExporterNoopTracing) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	return nil
}

func (e ExporterNoopTracing) Shutdown(ctx context.Context) error {
	return nil
}
