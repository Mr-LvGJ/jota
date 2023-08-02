package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	uhttp "github.com/Mr-LvGJ/http"
	"github.com/Mr-LvGJ/jota/access_log"
	"github.com/Mr-LvGJ/jota/log"
	"github.com/Mr-LvGJ/tracing"
	tracing_echo "github.com/Mr-LvGJ/tracing/http/echo"
)

func main() {
	// 创建一个 Echo 实例
	c := log.DefaultConfig()
	c.Filename = "server.log"
	if err := log.NewGlobal(c); err != nil {
		panic(err)
	}

	accessLogConfig := *c
	accessLogConfig.Filename = "server-access.log"

	e := echo.New()
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
	e.Use(uhttp.ReplaceWriterInterceptor())
	e.Use(middleware.RequestID())
	e.Use(tracing_echo.TracingEchoHTTPServerInterceptor())
	e.Use(access_log.AccessLogInterceptor(access_log.WithLogConfig(accessLogConfig)))
	// 定义路由和处理函数
	e.GET("/", helloHandler)
	e.GET("/echo/:id", echoHandler)
	// 启动服务器
	e.Start(":18080")
}

// 处理函数
func helloHandler(c echo.Context) error {
	log.Info(c.Request().Context(), "hello handler get", "header", c.Request().Header)
	return c.String(http.StatusOK, "Hello, World! Request ID: "+c.Response().Header().Get(echo.HeaderXRequestID))
}

func echoHandler(c echo.Context) error {
	log.Info(c.Request().Context(), "echo handler get", "header", c.Request().Header)
	return c.String(http.StatusOK, c.Param("id"))
}
