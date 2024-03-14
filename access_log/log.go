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
		var (
			bodyBytes []byte
			resp      any
		)
		requestId, _ := c.Get(XRequestIDKey)
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw
		attrs := []any{
			"request_id", requestId,
			"body", string(bodyBytes),
			"ip", c.ClientIP(),
			"raw_query", c.Request.URL.RawQuery,
			"user_agent", c.Request.UserAgent(),
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
		}

		accessLog.logger.Info(c, "handle request begin", attrs...)

		start := time.Now().UTC()
		c.Next()
		end := time.Now().UTC()

		if err := json.Unmarshal(blw.body.Bytes(), &resp); err != nil {
			attrs = append(attrs, "err", err)
		}
		attrs = append(attrs,
			"response", resp,
			"status_code", c.Writer.Status(),
			"cost", end.Sub(start),
		)
		if c.Writer.Status() >= 500 {
			accessLog.logger.Error(c, "handle request error", attrs...)
		} else {
			accessLog.logger.Info(c, "handle request done", attrs...)
		}
	}
}
