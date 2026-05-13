package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ElliAbby/go_cinema_system/internal/platform/kafka"
	"github.com/ElliAbby/go_cinema_system/internal/order"

)

func handlePaymentMessage(ctx context.Context, uc order.UseCase, payload []byte) error {
	var event kafka.PaymentRequestedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	start := time.Now()
	booking, tickets, err := uc.ProcessPayment(ctx, event.UserID, event.BookingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Booking %s already processed or not found", event.BookingID)
			return nil
		}

		log.Printf("Payment processing failed for booking %s: %v", event.BookingID, err)
		return nil
	}

	log.Printf("Booking %s paid successfully in %s, tickets=%d, total=%.2f",
		booking.ID, time.Since(start).String(), len(tickets), booking.TotalPrice)
	return nil
}
