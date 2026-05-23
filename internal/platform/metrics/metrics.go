package metrics

import (
	"log"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const businessMetricsRefreshInterval = 30 * time.Second

type BusinessMetrics struct {
	// HTTP метрики
	RequestDuration *prometheus.HistogramVec
	RequestCount    *prometheus.CounterVec

	// Метрики бизнес-логики
	ActiveBookings prometheus.Gauge
	TicketsSold prometheus.Gauge
	RevenueTotal   prometheus.Gauge

	// Дополнительные метрики
	ActiveUsers     prometheus.Gauge
	SessionsActive  prometheus.Gauge
	BookingErrors   *prometheus.CounterVec
	PaymentDuration prometheus.Histogram
}

var Metrics *BusinessMetrics

func InitMetrics(db *sqlx.DB, serviceName string) {
	registerer := promauto.With(labeledRegisterer(serviceName))

	Metrics = &BusinessMetrics{
		// HTTP метрики
		RequestDuration: registerer.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cinema_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),

		RequestCount: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "cinema_http_requests_total",
			Help: "Total number of HTTP requests.",
		}, []string{"method", "path", "status"}),

		// Бизнес метрики
		ActiveBookings: registerer.NewGauge(prometheus.GaugeOpts{
			Name: "cinema_active_bookings_total",
			Help: "Total number of active bookings (pending status)",
		}),

		TicketsSold: registerer.NewGauge(prometheus.GaugeOpts{
			Name: "cinema_tickets_sold_total",
			Help: "Total number of tickets sold (paid status)",
		}),

		RevenueTotal: registerer.NewGauge(prometheus.GaugeOpts{
			Name: "cinema_revenue_total_rub",
			Help: "Total revenue in rubles from paid bookings",
		}),

		// Дополнительные метрики
		ActiveUsers: registerer.NewGauge(prometheus.GaugeOpts{
			Name: "cinema_active_users_total",
			Help: "Total number of active users",
		}),

		SessionsActive: registerer.NewGauge(prometheus.GaugeOpts{
			Name: "cinema_active_sessions_total",
			Help: "Total number of active movie sessions",
		}),

		BookingErrors: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "cinema_booking_errors_total",
			Help: "Total number of booking errors by type",
		}, []string{"error_type"}),

		PaymentDuration: registerer.NewHistogram(prometheus.HistogramOpts{
			Name:    "cinema_payment_duration_seconds",
			Help:    "Duration of payment processing in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		}),
	}

	refreshBusinessMetrics(db)
	go runBusinessMetricsRefresh(db)
}

// увеличивает счетчик активных бронирований
func IncActiveBookings() {
	if Metrics != nil {
		Metrics.ActiveBookings.Inc()
	}
}

// уменьшает счетчик активных бронирований
func DecActiveBookings() {
	if Metrics != nil {
		Metrics.ActiveBookings.Dec()
	}
}

// увеличивает счетчик проданных билетов
func AddTicketsSold(count float64) {
	if Metrics != nil {
		Metrics.TicketsSold.Add(count)
	}
}

// устанавливает количество активных пользователей
func SetActiveUsers(count float64) {
	if Metrics != nil {
		Metrics.ActiveUsers.Set(count)
	}
}

// устанавливает количество активных сеансов
func SetSessionsActive(count float64) {
	if Metrics != nil {
		Metrics.SessionsActive.Set(count)
	}
}

// увеличивает общую выручку
func AddRevenue(amount float64) {
	if Metrics != nil {
		Metrics.RevenueTotal.Add(amount)
	}
}

// увеличивает счетчик ошибок бронирования по типу
func IncBookingErrors(errorType string) {
	if Metrics != nil {
		Metrics.BookingErrors.WithLabelValues(errorType).Inc()
	}
}

// записывает длительность обработки платежа
func ObservePaymentDuration(duration float64) {
	if Metrics != nil {
		Metrics.PaymentDuration.Observe(duration)
	}
}

func runBusinessMetricsRefresh(db *sqlx.DB) {
	ticker := time.NewTicker(businessMetricsRefreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		refreshBusinessMetrics(db)
	}
}

func refreshBusinessMetrics(db *sqlx.DB) {
	if db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var snapshot struct {
		ActiveBookings int64 `db:"active_bookings"`
		TicketsSold    int64 `db:"tickets_sold"`
		ActiveUsers    int64 `db:"active_users"`
		RevenueTotal   float64 `db:"revenue_total"`
		SessionsActive int64   `db:"sessions_active"`
	}

	query := `
		SELECT
			COALESCE((SELECT COUNT(*) FROM bookings WHERE status = 'pending'), 0) AS active_bookings,
			COALESCE((SELECT COUNT(*) FROM tickets), 0) AS tickets_sold,
			COALESCE((SELECT COUNT(*) FROM users WHERE is_active = TRUE), 0) AS active_users,
			COALESCE((SELECT SUM(total_price) FROM bookings WHERE status = 'paid'), 0) AS revenue_total,
			COALESCE((SELECT COUNT(*) FROM sessions), 0) AS sessions_active
	`
	if err := db.GetContext(ctx, &snapshot, query); err != nil {
		log.Printf("Failed to refresh business metrics: %v", err)
		return
	}

	if Metrics != nil {
		Metrics.ActiveBookings.Set(float64(snapshot.ActiveBookings))
		Metrics.TicketsSold.Set(float64(snapshot.TicketsSold))
		SetActiveUsers(float64(snapshot.ActiveUsers))
		Metrics.RevenueTotal.Set(snapshot.RevenueTotal)
		SetSessionsActive(float64(snapshot.SessionsActive))
	}
}