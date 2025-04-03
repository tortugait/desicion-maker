package http

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type wrappedResponseWriter struct {
	originalEchoResponseWriter http.ResponseWriter
	body                       []byte
}

func (w *wrappedResponseWriter) Header() http.Header {
	return w.originalEchoResponseWriter.Header()
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.originalEchoResponseWriter.Write(data)
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.originalEchoResponseWriter.WriteHeader(statusCode)
}

func (w *wrappedResponseWriter) Body() []byte {
	return w.body
}

func newLoggingMiddleware(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			startTime := time.Now()
			wrappedResponseWriter := &wrappedResponseWriter{
				originalEchoResponseWriter: eCtx.Response().Writer,
			}
			eCtx.Response().Writer = wrappedResponseWriter
			err := next(eCtx)
			duration := time.Since(startTime)

			status := eCtx.Response().Status
			params := []any{
				"method", eCtx.Request().Method,
				"url", eCtx.Request().URL.String(),
				"duration", duration.String(),
			}
			if err != nil {
				params = append(params, "error", err)

				var eErr *echo.HTTPError
				if errors.As(err, &eErr) {
					status = eErr.Code
				}
			}

			params = append(params, "status", status)

			logFn := logger.Debugw

			switch {
			case status >= http.StatusInternalServerError:
				logFn = logger.Errorw
			case duration >= time.Second*3:
				logFn = logger.Warnw
			default:
				requestBody := eCtx.Request().Body
				requestBodyBytes, _ := io.ReadAll(requestBody)
				params = append(params, "requestBody", string(requestBodyBytes),
					"responseBody", string(wrappedResponseWriter.Body()))
			}

			logFn("HTTP endpoint call", params...)

			return err
		}
	}
}
