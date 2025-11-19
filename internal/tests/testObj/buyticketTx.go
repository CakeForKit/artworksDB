package testobj

import (
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
)

type TicketPurchaseTxMother interface {
	TicketPurchaseTxP(
		tptxID uuid.UUID, eventID uuid.UUID, userID uuid.UUID, cntTickets int,
	) *models.TicketPurchaseTx
	TicketPurchaseTxByUserP(
		tptxID uuid.UUID, customerName string, customerEmail string,
		eventID uuid.UUID, userID uuid.UUID, cntTickets int,
	) *models.TicketPurchaseTx
}

func NewTicketPurchaseTxMother() TicketPurchaseTxMother {
	return &tptxMother{}
}

type tptxMother struct {
}

func (um *tptxMother) TicketPurchaseTxP(
	tptxID uuid.UUID, eventID uuid.UUID, userID uuid.UUID, cntTickets int) *models.TicketPurchaseTx {
	tptx, _ := models.NewBuyTicketTx(
		tptxID,
		"test-customer-name",
		"test@mail.ru",
		time.Now(),
		eventID,
		userID,
		cntTickets,
		time.Now().Add(time.Hour),
	)
	return &tptx
}

func (um *tptxMother) TicketPurchaseTxByUserP(tptxID uuid.UUID, customerName string, customerEmail string,
	eventID uuid.UUID, userID uuid.UUID, cntTickets int) *models.TicketPurchaseTx {
	tptx, _ := models.NewBuyTicketTx(
		tptxID,
		customerName,
		customerEmail,
		time.Now(),
		eventID,
		userID,
		cntTickets,
		time.Now().Add(time.Hour),
	)
	return &tptx
}
