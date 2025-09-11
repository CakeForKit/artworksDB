package testobj

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type TicketPurchaseMother interface {
	TicketPurchaseP(tptxID uuid.UUID, eventID uuid.UUID, userID uuid.UUID) *models.TicketPurchase
}

func NewTicketPurchaseMother() TicketPurchaseMother {
	return &ticketPurchaseMother{}
}

type ticketPurchaseMother struct {
}

func (um *ticketPurchaseMother) TicketPurchaseP(tpID uuid.UUID, eventID uuid.UUID, userID uuid.UUID) *models.TicketPurchase {
	tp, _ := models.NewTicketPurchase(
		tpID,
		"test-customer-name",
		"test-customer-email",
		time.Now(),
		eventID,
		userID,
	)
	return &tp
}
