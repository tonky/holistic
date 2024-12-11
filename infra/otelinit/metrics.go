package otelinit

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type ExporterType string

const (
	Noop   ExporterType = "noop"
	Stdout ExporterType = "stdout"
	Zipkin ExporterType = "zipkin"
	Jaeger ExporterType = "jaeger"
)

type Config struct {
	RegisterGlobal bool         `default:"true"`
	ExporterType   ExporterType `default:"noop"`
}

var Meter metric.Meter

func InitMetrics(conf Config, serviceName string) (metric.Meter, func(), error) {
	// Create resource.
	res, err := newResource(serviceName)
	if err != nil {
		return nil, nil, err
	}

	// Create a meter provider.
	// You can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.
	meterProvider, err := newMeterProvider(res, conf.ExporterType)
	if err != nil {
		return nil, nil, err
	}

	// Register as global meter provider so that it can be used via otel.Meter
	// and accessed using otel.GetMeterProvider.
	// Most instrumentation libraries use the global meter provider as default.
	// If the global meter provider is not set then a no-op implementation
	// is used, which fails to generate data.
	otel.SetMeterProvider(meterProvider)

	if conf.RegisterGlobal {
		Meter = meterProvider.Meter(serviceName)
	}

	// Handle shutdown properly so nothing leaks.
	return meterProvider.Meter(serviceName), func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}, nil
}

func newResource(serviceName string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("0.1.0"),
		))
}

func newMeterProvider(res *resource.Resource, et ExporterType) (*sdkmetric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	if et == Noop {
		metricExporter = ExporterNoopMeter{}
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			sdkmetric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}

type ExporterNoopMetric struct{}

func (mp ExporterNoopMetric) ForceFlush(ctx context.Context) error { return nil }
func (mp ExporterNoopMetric) Meter(name string, options ...metric.MeterOption) metric.Meter {
	return nil
}

type ExporterNoopMeter struct{}

func (e ExporterNoopMeter) Temporality(sdkmetric.InstrumentKind) metricdata.Temporality { return 0 }
func (e ExporterNoopMeter) Aggregation(sdkmetric.InstrumentKind) sdkmetric.Aggregation  { return nil }
func (e ExporterNoopMeter) Export(context.Context, *metricdata.ResourceMetrics) error   { return nil }
func (e ExporterNoopMeter) ForceFlush(context.Context) error                            { return nil }
func (e ExporterNoopMeter) Shutdown(context.Context) error                              { return nil }
