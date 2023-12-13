package access_log

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

const XRequestIDKey = "X-Request-ID"

func AccessLogInterceptor(opts ...option) echo.MiddlewareFunc {
	newConfig(opts...)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			ctx := c.Request().Context()

			err := next(c)

			if err != nil {
				accessLog.logger.Error(ctx, "Handle request done", "err", err)
			} else {
				accessLog.logger.Warn(ctx, "Handle request done",
					"Method", c.Request().Method,
					"Path", c.Path(),
					"__source__", c.RealIP(),
					"Cost", time.Since(startTime))
			}
			return err
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinAccessLogInterceptor(opts ...option) gin.HandlerFunc {
	newConfig(opts...)
	return func(c *gin.Context) {
		requestId, _ := c.Get(XRequestIDKey)
		start := time.Now().UTC()
		path := c.Request.URL.Path
		var (
			bodyBytes []byte
			resp      any
		)

		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		method := c.Request.Method
		ip := c.ClientIP()
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw
		accessLog.logger.Info(c, "handle request begin",
			"requestId", requestId,
			"body", string(bodyBytes),
			"ip", ip,
			"path", path,
			"method", method)
		c.Next()
		end := time.Now().UTC()
		latency := end.Sub(start)
		if err := json.Unmarshal(blw.body.Bytes(), &resp); err != nil {
			accessLog.logger.Error(c, "handle request done",
				"requestId", requestId,
				"latency", latency,
				"ip", ip,
				"path", path,
				"method", method,
				"err", err)
		} else {
			accessLog.logger.Info(c, "handle request done",
				"requestId", requestId,
				"response", resp,
				"latency", latency,
				"ip", ip,
				"path", path,
				"method", method)
		}

	}
}
