package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"

)

type BusinessMetrics struct {
	// HTTP метрики
	RequestDuration *prometheus.HistogramVec
	RequestCount *prometheus.CounterVec

	// Метрики бизнес-логики
	ActiveBookings prometheus.Gauge
	TicketsSold prometheus.Counter
	RevenueTotal prometheus.Counter

	// Дополнительные метрики
	ActiveUsers prometheus.Gauge
	SessionsActive prometheus.Gauge
	BookingErrors *prometheus.CounterVec
	PaymentDuration prometheus.Histogram
}

var Metrics *BusinessMetrics

func InitMetrics() {
	Metrics = &BusinessMetrics{
		// HTTP метрики
		RequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cinema_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),

		RequestCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "cinema_http_requests_total",
			Help: "Total number of HTTP requests.",
        }, []string{"method", "path", "status"}),
        
        // Бизнес метрики
        ActiveBookings: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "cinema_active_bookings_total",
            Help: "Total number of active bookings (pending status)",
        }),
        
        TicketsSold: promauto.NewCounter(prometheus.CounterOpts{
            Name: "cinema_tickets_sold_total",
            Help: "Total number of tickets sold (paid status)",
        }),
        
        RevenueTotal: promauto.NewCounter(prometheus.CounterOpts{
            Name: "cinema_revenue_total_rub",
            Help: "Total revenue in rubles from paid bookings",
        }),
        
        // Дополнительные метрики
        ActiveUsers: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "cinema_active_users_total",
            Help: "Total number of active users",
        }),
        
        SessionsActive: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "cinema_active_sessions_total",
            Help: "Total number of active movie sessions",
        }),
        
        BookingErrors: promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "cinema_booking_errors_total",
            Help: "Total number of booking errors by type",
        }, []string{"error_type"}),
        
        PaymentDuration: promauto.NewHistogram(prometheus.HistogramOpts{
            Name:    "cinema_payment_duration_seconds",
            Help:    "Duration of payment processing in seconds",
            Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
        }),
    }
}

func IncActiveBookings() {
	if Metrics != nil {
		Metrics.ActiveBookings.Inc()
	}
}

func DecActiveBookings() {
	if Metrics != nil {
		Metrics.ActiveBookings.Dec()
	}
}

func AddTicketsSold(count float64) {
	if Metrics != nil {
		Metrics.TicketsSold.Add(count)
	}
}

func SetActiveUsers(count float64) {
	if Metrics != nil {
		Metrics.ActiveUsers.Set(count)
	}
}

func SetSessionsActive(count float64) {
	if Metrics != nil {
		Metrics.SessionsActive.Set(count)
	}
}

func AddRevenue(amount float64) {
	if Metrics != nil {
		Metrics.RevenueTotal.Add(amount)
	}
}

func IncBookingErrors(errorType string) {
	if Metrics != nil {
		Metrics.BookingErrors.WithLabelValues(errorType).Inc()
	}
}

func ObservePaymentDuration(duration float64) {
	if Metrics != nil {
		Metrics.PaymentDuration.Observe(duration)
	}
}
