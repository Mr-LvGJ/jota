package echo

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const defaultName = "github.com/Mr-LvGJ/tracing/http/echo"

type config struct {
	tracer            trace.Tracer
	propagators       propagation.TextMapPropagator
	spanNameFormatter func(string, *http.Request) string
	tracerProvider    trace.TracerProvider
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func newTracer(tp trace.TracerProvider) trace.Tracer {
	return tp.Tracer(defaultName)
}

func newConfig(opts ...Option) *config {
	c := &config{
		propagators: otel.GetTextMapPropagator(),
	}
	for _, opt := range opts {
		opt.apply(c)
	}

	if c.spanNameFormatter == nil {
		c.spanNameFormatter = func(s string, r *http.Request) string {
			return "HTTP " + s + " " + r.Method + " " + r.URL.Path
		}
	}

	return c
}
