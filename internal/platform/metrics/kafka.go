package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type KafkaMetrics struct {
	// Consumer метрики (получение сообщений)
	MessagesReceived              *prometheus.CounterVec // сколько сообщений получено
	MessagesProcessed             *prometheus.CounterVec // сколько успешно обработано
	MessagesFailed                *prometheus.CounterVec // сколько сообщений не удалось обработать
	MessageRetries                *prometheus.CounterVec // сколько сообщений было повторно отправлено
	ConsumerErrors                *prometheus.CounterVec // сколько ошибок произошло в консьюмере
	MessageProcessingDuration     *prometheus.HistogramVec // сколько времени ушло на обработку сообщения
	// Producer метрики (отправка сообщений)
	MessagesProduced              *prometheus.CounterVec // сколько сообщений было отправлено
	ProducerDeliveryDuration      *prometheus.HistogramVec // сколько времени ушло на доставку сообщения
	ProducerErrors                *prometheus.CounterVec // сколько ошибок произошло в продюсере
	
	MessagesSkipped               *prometheus.CounterVec // сколько сообщений было пропущено
}

var KafkaStats *KafkaMetrics

func InitKafkaMetrics(serviceName string, topic string) {
	registerer := promauto.With(labeledRegisterer(serviceName))

	KafkaStats = &KafkaMetrics{
		MessagesReceived: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "kafka_messages_received_total",
			Help: "Total number of Kafka messages received by the consumer.",
		}, []string{"topic"}),
		MessagesProcessed: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "kafka_messages_processed_total",
			Help: "Total number of Kafka messages processed successfully.",
		}, []string{"topic"}),
		MessagesFailed: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "kafka_messages_failed_total",
			Help: "Total number of Kafka messages that failed processing.",
		}, []string{"topic", "error_type"}),
		MessageRetries: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "kafka_message_retries_total",
			Help: "Total number of Kafka message retries.",
		}, []string{"topic"}),
		ConsumerErrors: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "kafka_consumer_errors_total",
			Help: "Total number of Kafka consumer errors by type.",
		}, []string{"topic", "error_type"}),
		MessageProcessingDuration: registerer.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "kafka_message_processing_duration_seconds",
			Help:    "Duration of Kafka message processing in seconds.",
			Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10},
		}, []string{"topic"}),
		MessagesProduced: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "kafka_messages_produced_total",
			Help: "Total number of Kafka messages produced.",
		}, []string{"topic"}),
		ProducerDeliveryDuration: registerer.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "kafka_producer_delivery_duration_seconds",
			Help:    "Duration of Kafka message delivery in seconds.",
			Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
		}, []string{"topic"}),
		ProducerErrors: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "kafka_producer_errors_total",
			Help: "Total number of Kafka producer errors by type.",
		}, []string{"topic", "error_type"}),
		MessagesSkipped: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "messages_skipped_total",
			Help: "Total number of skipped Kafka messages by service.",
		}, []string{"service"}),
	}

	if topic != "" {
		KafkaStats.MessagesReceived.WithLabelValues(topic)
		KafkaStats.MessagesProcessed.WithLabelValues(topic)
		KafkaStats.MessageRetries.WithLabelValues(topic)
		KafkaStats.MessageProcessingDuration.WithLabelValues(topic)
		KafkaStats.MessagesProduced.WithLabelValues(topic)
		KafkaStats.ProducerDeliveryDuration.WithLabelValues(topic)
	}
	if topic != "" {
		KafkaStats.MessagesFailed.WithLabelValues(topic, "unknown")
		KafkaStats.ConsumerErrors.WithLabelValues(topic, "unknown")
		KafkaStats.ProducerErrors.WithLabelValues(topic, "unknown")
	}
	if serviceName != "" {
		KafkaStats.MessagesSkipped.WithLabelValues(serviceName)
	}
}
 
func IncKafkaMessagesReceived(topic string) {
	if KafkaStats != nil {
		KafkaStats.MessagesReceived.WithLabelValues(topic).Inc()
	}
}

func IncKafkaMessagesProcessed(topic string) {
	if KafkaStats != nil {
		KafkaStats.MessagesProcessed.WithLabelValues(topic).Inc()
	}
}

func IncKafkaMessagesFailed(topic string, errorType string) {
	if KafkaStats != nil {
		KafkaStats.MessagesFailed.WithLabelValues(topic, errorType).Inc()
	}
}

func IncKafkaMessageRetries(topic string) {
	if KafkaStats != nil {
		KafkaStats.MessageRetries.WithLabelValues(topic).Inc()
	}
}

func IncKafkaConsumerErrors(topic string, errorType string) {
	if KafkaStats != nil {
		KafkaStats.ConsumerErrors.WithLabelValues(topic, errorType).Inc()
	}
}

func ObserveKafkaMessageProcessingDuration(topic string, duration float64) {
	if KafkaStats != nil {
		KafkaStats.MessageProcessingDuration.WithLabelValues(topic).Observe(duration)
	}
}

func IncKafkaMessagesProduced(topic string) {
	if KafkaStats != nil {
		KafkaStats.MessagesProduced.WithLabelValues(topic).Inc()
	}
}

func ObserveKafkaProducerDeliveryDuration(topic string, duration float64) {
	if KafkaStats != nil {
		KafkaStats.ProducerDeliveryDuration.WithLabelValues(topic).Observe(duration)
	}
}

func IncKafkaProducerErrors(topic string, errorType string) {
	if KafkaStats != nil {
		KafkaStats.ProducerErrors.WithLabelValues(topic, errorType).Inc()
	}
}

func IncMessagesSkipped(service string) {
	if KafkaStats != nil {
		KafkaStats.MessagesSkipped.WithLabelValues(service).Inc()
	}
}
