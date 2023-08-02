package access_log

import (
	"time"

	"github.com/labstack/echo/v4"
)

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
