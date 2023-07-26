package echo

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/Mr-LvGJ/http"
	"github.com/Mr-LvGJ/log"
)

const TraceID = "trace-id"

func TracingEchoHTTPServerInterceptor(opts ...Option) echo.MiddlewareFunc {
	config := newConfig(opts...)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tracer := config.tracer
			if tracer == nil {
				if span := trace.SpanFromContext(c.Request().Context()); span.SpanContext().IsValid() {
					tracer = newTracer(span.TracerProvider())
				} else {
					tracer = newTracer(otel.GetTracerProvider())
				}
			}

			ctx := config.propagators.Extract(c.Request().Context(), propagation.HeaderCarrier(c.Request().Header))
			ctx, span := tracer.Start(ctx, config.spanNameFormatter("Server", c.Request()))
			defer span.End()

			ctx = log.WithValue(ctx, TraceID, span.SpanContext().TraceID())

			c.SetRequest(c.Request().WithContext(ctx))
			span.SetAttributes(attribute.Key(TraceID).String(span.SpanContext().TraceID().String()))
			span.SetAttributes(attribute.String(echo.HeaderXRequestID, c.Response().Header().Get(echo.HeaderXRequestID)))
			err := next(c)
			//if c.Response().Status != http.StatusOK {
			span.SetAttributes(attribute.String("code", strconv.Itoa(c.Response().Status)))
			span.SetAttributes(attribute.String("response", string(c.Response().Writer.(*http.ResponseBuffer).Content())))
			span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(c.Response().Status))
			//}
			return err
		}
	}
}
