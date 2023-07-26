package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	exporter_jaeger "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	DefaultServiceName        = "default"
	DefaultServiceVersion     = ""
	DefaultExporterEndpoint   = "127.0.0.1:4317"
	DefaultMaxQueueSize       = 2048
	DefaultMaxExportBatchSize = 512
)

var conf = &Config{
	ServiceName:                   DefaultServiceName,
	ServiceVersion:                DefaultServiceVersion,
	TraceExporterEndpoint:         DefaultExporterEndpoint,
	TraceExporterEndpointInsecure: true,

	MaxQueueSize:       DefaultMaxQueueSize,
	MaxExportBatchSize: DefaultMaxExportBatchSize,
}

type Config struct {
	ServiceName                   string
	ServiceVersion                string // using xversion
	ServiceCommitId               string
	ServiceBuildTime              string
	TraceExporterEndpoint         string
	TraceExporterEndpointInsecure bool
	MaxQueueSize                  int
	MaxExportBatchSize            int
}

type Option func(*Config)

func WithTraceExporterEndpoint(agentAddr string) Option {
	return func(c *Config) {
		c.TraceExporterEndpoint = agentAddr
	}
}

func WithServiceName(serviceName string) Option {
	return func(c *Config) {
		c.ServiceName = serviceName
	}
}

func WithServiceVersion(serviceVersion string) Option {
	return func(c *Config) {
		c.ServiceVersion = serviceVersion
	}
}

func InitTracer(options ...Option) *sdktrace.TracerProvider {
	for _, option := range options {
		option(conf)
	}

	exporter, err := initJaegerExporter()
	if err != nil {
		// panic when init phase failed
		panic(err)
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(conf.MaxExportBatchSize),
			sdktrace.WithMaxQueueSize(conf.MaxQueueSize)),
		sdktrace.WithResource(getDefaultResource()),
	}

	tp := sdktrace.NewTracerProvider(opts...)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
		jaeger.Jaeger{}),
	)

	return tp
}

func initExporter() (*otlptrace.Exporter, error) {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(conf.TraceExporterEndpoint),
	}

	if conf.TraceExporterEndpointInsecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}
	exporter, err := otlptracegrpc.New(context.Background(), options...)

	if err != nil {
		panic(err)
	}

	return exporter, nil
}

func initJaegerExporter() (*exporter_jaeger.Exporter, error) {
	options := []exporter_jaeger.CollectorEndpointOption{
		exporter_jaeger.WithEndpoint(conf.TraceExporterEndpoint),
	}
	exporter, err := exporter_jaeger.New(exporter_jaeger.WithCollectorEndpoint(options...))
	if err != nil {
		panic(err)
	}
	return exporter, nil
}

func getDefaultResource() *resource.Resource {
	hostname, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(conf.ServiceName),
		semconv.HostNameKey.String(hostname),
		semconv.ServiceVersionKey.String(conf.ServiceVersion),
		semconv.ProcessPIDKey.Int(os.Getpid()),
		semconv.ProcessCommandKey.String(os.Args[0]),
	)
}
