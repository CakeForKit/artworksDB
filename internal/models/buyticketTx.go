package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type BuyTicketTx struct {
	id         uuid.UUID
	userID     uuid.UUID
	eventID    uuid.UUID
	cntTickets int
	expiredAt  time.Time
}

var (
	ErrBuyTicketTxZeroCnt = errors.New("cntTickets <= 0")
)

func NewBuyTicketTx(id uuid.UUID, userID uuid.UUID, eventID uuid.UUID, cntTickets int, expiredAt time.Time) (BuyTicketTx, error) {
	if cntTickets <= 0 {
		return BuyTicketTx{}, ErrBuyTicketTxZeroCnt
	}
	return BuyTicketTx{
		id:         id,
		userID:     userID,
		eventID:    eventID,
		cntTickets: cntTickets,
		expiredAt:  expiredAt,
	}, nil
}

func (t *BuyTicketTx) GetID() uuid.UUID {
	return t.id
}
