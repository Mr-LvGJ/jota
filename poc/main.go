package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Mr-LvGJ/log"

	"github.com/Mr-LvGJ/tracing"
)

type ctxKey int

const contextLogKey ctxKey = iota

func main() {
	ctx := context.Background()
	logCfg := log.DefaultConfig()
	logCfg.Filename = "./poc.log"
	//logCfg.Prod = true
	log.NewGlobal(logCfg)
	ac := "sss"
	aMap := map[string]interface{}{
		"acc": "ccc",
		"bbb": 1,
	}
	tp := tracing.InitTracer(
		tracing.WithTraceExporterEndpoint("http://devbox:14268/api/traces"),
		tracing.WithServiceName("poc"),
	)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Error(ctx, "Error shutting down tracer provider:", "err", err)
		}
	}()

	ctx, span := StartFromContext(ctx, "cc", "test")
	ctx = log.WithValue(ctx, "trace-id", span.SpanContext().TraceID())
	defer span.End()
	span.SetAttributes(
		attribute.String("query.username", "c"),
	)
	time.Sleep(3 * time.Second)
	span.SetName("cccccccccc")
	_, span2 := StartFromContext(ctx, "cc2", "function 2")
	span2.SetName("c222")
	span.SetAttributes(attribute.String("funcation2", "true"))
	_, span3 := StartFromContext(ctx, "windu", "ListResourcePermission")
	span3.SetName("Windu")
	defer span2.End()
	defer span3.End()
	log.Error(ctx, "ccc", "ac", ac, "aMap", aMap)
	log.Info(ctx, "ctx", "ccc", "cccc")
	log.Info(ctx, "ctx", "cc", ctx.Value(contextLogKey))
}

func StartFromContext(ctx context.Context, tracer, spanName string) (context.Context, trace.Span) {
	tp := otel.GetTracerProvider()
	t := tp.Tracer(tracer)
	return t.Start(ctx, spanName)
}
