package order

import "errors"

var ErrReservationExpired = errors.New("reservation expired: payment time limit exceeded")