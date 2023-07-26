package promonitor

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"strconv"
)

var defaultMetricPath = "/metrics"
var defaultConfig = &config{
	metricPath: defaultMetricPath,
}

var reqCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "gin_prometheus",
		Name:      "http_request_count",
		Help:      "Total number of request count",
	},
	[]string{"http_method", "status"},
)

type config struct {
	listenAddress string
	metricPath    string
}
type Option func(*config)

func WithListenAddress(listenAddress string) func(*config) {
	return func(c *config) {
		c.listenAddress = listenAddress
	}
}

func WithMetricPath(metricPath string) func(*config) {
	return func(c *config) {
		c.metricPath = metricPath
	}
}

func newConfig(opts ...Option) {
	for _, opt := range opts {
		opt(defaultConfig)
	}
}

func NewGinPrometheusMiddleware(opts ...Option) gin.HandlerFunc {
	newConfig(opts...)
	return func(c *gin.Context) {
		if c.Request.URL.Path == defaultConfig.metricPath ||
			c.Request.URL.Path == "/favicon.ico" {
			c.Next()
			return
		}
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		reqCount.WithLabelValues(c.Request.URL.Path, status)
	}
}

func init() {
	prometheus.MustRegister(reqCount)
	prometheus.Register(collectors.NewGoCollector())
	prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}
