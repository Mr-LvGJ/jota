package http

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseBuffer struct {
	writer http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *ResponseBuffer) Header() http.Header {
	return r.writer.Header()
}

func (r *ResponseBuffer) Write(body []byte) (int, error) {
	r.body.Write(body)
	return r.writer.Write(body)
}

func (r *ResponseBuffer) WriteHeader(statusCode int) {
	r.status = statusCode
	r.writer.WriteHeader(statusCode)
}

func (r *ResponseBuffer) Content() []byte {
	return r.body.Bytes()
}

func ReplaceWriterInterceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rec := &ResponseBuffer{writer: c.Response().Writer}
			c.Response().Writer = rec
			return next(c)
		}
	}
}
