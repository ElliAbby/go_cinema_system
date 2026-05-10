package http

import (
	"net/http"
	"strconv"
	"time"

	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ElliAbby/go_cinema_system/internal/metrics"

)

// var (
// 	// Время выполнения HTTP запросов в секундах
// 	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
// 		Name:    "cinema_http_request_duration_seconds",
// 		Help:    "Duration of HTTP requests in seconds.",
// 		Buckets: prometheus.DefBuckets,
// 	}, []string{"method", "path", "status"})
	
// 	// Общее количество HTTP запросов
// 	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
// 		Name: "cinema_http_requests_total",
// 		Help: "Total number of HTTP requests.",
// 	}, []string{"method", "path", "status"})

// 	// Активные соединения
//     activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
//         Name: "cinema_active_connections",
//         Help: "Number of active HTTP connections",
//     })
    
//     // Размер запросов
//     requestSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
//         Name:    "cinema_http_request_size_bytes",
//         Help:    "Size of HTTP requests in bytes",
//         Buckets: []float64{100, 500, 1000, 5000, 10000, 50000},
//     }, []string{"method", "path"})

// 	activeBookings = promauto.NewGauge(prometheus.GaugeOpts{
//         Name: "cinema_active_bookings_total",
//         Help: "Total number of active bookings",
//     })
    
//     ticketsSold = promauto.NewCounter(prometheus.CounterOpts{
//         Name: "cinema_tickets_sold_total",
//         Help: "Total number of tickets sold",
//     })
    
//     revenueTotal = promauto.NewCounter(prometheus.CounterOpts{
//         Name: "cinema_revenue_total_rub",
//         Help: "Total revenue in rubles",
//     })
// )

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
		// activeConnections.Inc()
		// defer activeConnections.Dec()

		// if r.ContentLength > 0 {
		// 	requestSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(r.ContentLength))
		// }

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