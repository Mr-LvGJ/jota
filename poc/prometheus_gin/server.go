package main

import (
	prom_gin "github.com/Mr-LvGJ/promonitor/gin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	router := gin.Default()

	router.Use(prom_gin.NewGinPrometheusMiddleware(
		prom_gin.WithListenAddress(":10086"),
		prom_gin.WithMetricPath("/metrics"),
	))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.Run(":8080")
}
