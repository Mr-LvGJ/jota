package main

import (
	"net/http"

	"github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	p := echoprometheus.MetricsMiddleware()
	e.Use(p)
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Start(":8080")
}
