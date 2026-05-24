package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ElliAbby/go_cinema_system/internal/order"
	"github.com/ElliAbby/go_cinema_system/internal/platform/metrics"
	"github.com/ElliAbby/go_cinema_system/internal/platform/kafka"

)

func handlePaymentMessage(ctx context.Context, topic string, serviceName string, uc order.UseCase, payload []byte) error {
	start := time.Now()
	defer func() {
		metrics.ObserveKafkaMessageProcessingDuration(topic, time.Since(start).Seconds())
	}()

	var event kafka.PaymentRequestedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		metrics.IncKafkaConsumerErrors(topic, "decode_error")
		return err
	}

	booking, tickets, err := uc.ProcessPayment(ctx, event.UserID, event.BookingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Бронирование %s уже обработано или не найдено", event.BookingID)
			metrics.IncMessagesSkipped(serviceName)
			return nil
		}

		log.Printf("Обработка оплаты для бронирования %s не удалась: %v", event.BookingID, err)
		metrics.IncKafkaConsumerErrors(topic, "processing_error")
		return nil
	}

	log.Printf("Бронирование %s успешно оплачено за %s, билетов=%d, сумма=%.2f",
		booking.ID, time.Since(start).String(), len(tickets), booking.TotalPrice)

	metrics.ObservePaymentDuration(time.Since(start).Seconds())
	return nil
}
