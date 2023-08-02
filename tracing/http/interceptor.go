package http

import (
	ht "net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/Mr-LvGJ/jota/log"

	"github.com/Mr-LvGJ/http"
)

const TraceID = "trace-id"

func TracingHTTPClientInterceptor(opts ...Option) http.ClientInterceptor {
	c := newConfig(opts...)
	return func(req *ht.Request, requester http.Requester) (*ht.Response, error) {
		tracer := c.tracer
		if tracer == nil {
			if span := trace.SpanFromContext(req.Context()); span.SpanContext().IsValid() {
				tracer = newTracer(span.TracerProvider())
			} else {
				tracer = newTracer(otel.GetTracerProvider())
			}
		}

		ctx, span := tracer.Start(req.Context(), c.spanNameFormatter("Client", req))
		defer span.End()

		span.SetAttributes(attribute.Key(TraceID).String(span.SpanContext().TraceID().String()))

		c.propagators.Inject(ctx, propagation.HeaderCarrier(req.Header))
		req = req.WithContext(log.WithValue(ctx, TraceID, span.SpanContext().TraceID()))

		res, err := requester(req)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return res, err
	}
}
