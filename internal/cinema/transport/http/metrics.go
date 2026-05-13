package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ElliAbby/go_cinema_system/internal/platform/metrics"

)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		path := r.URL.Path
		status := strconv.Itoa(wrappedWriter.statusCode)
		duration := time.Since(startedAt).Seconds()

		if metrics.Metrics != nil {
			metrics.Metrics.RequestDuration.WithLabelValues(r.Method, path, status).Observe(duration)
			metrics.Metrics.RequestCount.WithLabelValues(r.Method, path, status).Inc()
		}
	})
}

func metricsHandler() http.Handler {
	return promhttp.Handler()
}