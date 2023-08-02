package main

import (
	"context"
	"io"
	"net/http"
	"time"

	m_http "github.com/Mr-LvGJ/http"
	"github.com/Mr-LvGJ/jota/log"
	"github.com/Mr-LvGJ/tracing"
	tracing_http "github.com/Mr-LvGJ/tracing/http"

	. "github.com/Mr-LvGJ/poc/init"
)

func main() {
	ctx := context.Background()
	ctx = log.WithValue(ctx, "test-log", "1999999999999")
	InitLog("client.log")
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

	//req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:18080/", nil)
	req, err := http.NewRequest(http.MethodGet, "http://localhost:18080/", nil)
	if err != nil {
		return
	}
	client := m_http.NewClient(m_http.WithClientInterceptors(
		tracing_http.TracingHTTPClientInterceptor(),
		HTTPClientHelloInterceptor(),
	))

	response, err := client.Do(req)
	if err != nil {
		log.Error(req.Context(), "client do error", "err", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(req.Context(), "io read all error", "err", err)
		return
	}
	log.Info(req.Context(), "response", "body", string(body))
	log.Info(response.Request.Context(), "response background context", "body", string(body))
}

func HTTPClientHelloInterceptor() m_http.ClientInterceptor {
	return func(req *http.Request, requester m_http.Requester) (*http.Response, error) {
		req.Header.Add("test", "HTTPClientHelloInterceptor")
		log.Info(req.Context(), "HTTPClientHelloInterceptor before")
		resp, err := requester(req)
		log.Info(req.Context(), "HTTPClientHelloInterceptor after", "err", err)
		return resp, err
	}
}
