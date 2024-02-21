package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(level slog.Level) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			request := slog.Group("request",
				slog.String("id", req.Header.Get("X-Request-ID")),
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
			)

			rw := &statusRecorderResponseWriter{
				ResponseWriter: w,
			}

			startTime := time.Now()
			next.ServeHTTP(rw, req)
			duration := time.Since(startTime)

			response := slog.Group("response",
				slog.Duration("runtime", duration),
				slog.Int("status", rw.statusCode),
			)

			slog.Default().Log(req.Context(), level, "HTTP Request", request, response)
		})
	}
}

type statusRecorderResponseWriter struct {
	http.ResponseWriter

	headerWritten bool
	statusCode    int
}

func (w *statusRecorderResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *statusRecorderResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)

	if !w.headerWritten {
		w.statusCode = statusCode
		w.headerWritten = true
	}
}

func (w *statusRecorderResponseWriter) Write(b []byte) (int, error) {
	if !w.headerWritten {
		w.statusCode = http.StatusOK
	}
	w.headerWritten = true

	return w.ResponseWriter.Write(b)
}
