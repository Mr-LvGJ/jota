package promonitor

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	[]string{"http_method", "path", "status"},
)

var reqHandledCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "gin_prometheus",
		Name:      "http_server_handled_total",
		Help:      "Total number of request completed",
	},
	[]string{"http_method", "path", "status"},
)

var reqDur = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "gin_prometheus",
		Name:      "http_server_handling_seconds",
		Help:      "Time spend of request",
	}, []string{"http_method", "path"},
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
		method := c.Request.Method
		path := c.Request.URL.Path
		if path == defaultConfig.metricPath ||
			path == "/favicon.ico" {
			c.Next()
			return
		}

		timer := prometheus.NewTimer(reqDur.WithLabelValues(method, path))
		c.Next()
		timer.ObserveDuration()

		status := strconv.Itoa(c.Writer.Status())
		reqCount.WithLabelValues(method, path, status).Inc()
		reqHandledCount.WithLabelValues(method, path, status).Inc()
	}
}

func init() {
	prometheus.MustRegister(reqCount, reqHandledCount, reqDur)
	prometheus.Register(collectors.NewGoCollector())
	prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}
