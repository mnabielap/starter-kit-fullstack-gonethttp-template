package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Logger(next http.Handler) http.Handler {
	// simple logger setup
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)

		logger.Info("Request Processed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrappedWriter.status),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", r.RemoteAddr),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}