package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type TicketPurchaseTx struct {
	ticketPurchase TicketPurchase
	cntTickets     int
	expiredAt      time.Time
}

type jsonTicketPurchaseTx struct {
	TicketPurchase jsonTicketPurchase
	CntTickets     int       `json:"cntTickets"`
	ExpiredAt      time.Time `json:"expiredAt"`
}

var (
	ErrValidateTicketTx   = errors.New("invalid model TicketTx")
	ErrBuyTicketTxZeroCnt = errors.New("cntTickets <= 0")
)

func NewBuyTicketTx(
	id uuid.UUID,
	customerName string,
	customerEmail string,
	purchaseDate time.Time,
	eventID uuid.UUID,
	userID uuid.UUID,
	cntTickets int,
	expiredAt time.Time,
) (TicketPurchaseTx, error) {
	if cntTickets <= 0 {
		return TicketPurchaseTx{}, fmt.Errorf("%w: %v", ErrValidateTicketTx, ErrBuyTicketTxZeroCnt)
	}
	tp, err := NewTicketPurchase(id, customerName, customerEmail, purchaseDate, eventID, userID)
	if err != nil {
		return TicketPurchaseTx{}, fmt.Errorf("%w: %v", ErrValidateTicketTx, err)
	}
	return TicketPurchaseTx{
		ticketPurchase: tp,
		cntTickets:     cntTickets,
		expiredAt:      expiredAt,
	}, nil
}

func (t *TicketPurchaseTx) Tojson() ([]byte, error) {
	ticketPurchaseJson := jsonTicketPurchase{
		ID:            t.ticketPurchase.id,
		CustomerName:  t.ticketPurchase.customerName,
		CustomerEmail: t.ticketPurchase.customerEmail,
		PurchaseDate:  t.ticketPurchase.purchaseDate,
		EventID:       t.ticketPurchase.eventID,
		UserID:        t.ticketPurchase.userID,
	}

	txJson := jsonTicketPurchaseTx{
		TicketPurchase: ticketPurchaseJson,
		CntTickets:     t.cntTickets,
		ExpiredAt:      t.expiredAt,
	}

	return json.Marshal(txJson)
}

func (t *TicketPurchaseTx) FromJson(data []byte) error {
	var txJson jsonTicketPurchaseTx
	if err := json.Unmarshal(data, &txJson); err != nil {
		return err
	}

	t.cntTickets = txJson.CntTickets
	t.expiredAt = txJson.ExpiredAt
	t.ticketPurchase = TicketPurchase{
		id:            txJson.TicketPurchase.ID,
		customerName:  txJson.TicketPurchase.CustomerName,
		customerEmail: txJson.TicketPurchase.CustomerEmail,
		purchaseDate:  txJson.TicketPurchase.PurchaseDate,
		eventID:       txJson.TicketPurchase.EventID,
		userID:        txJson.TicketPurchase.UserID,
	}

	return nil
}

func (t *TicketPurchaseTx) ToTxTicketPurchaseResponse() jsonreqresp.TxTicketPurchaseResponse {
	ticketPurchaseJson := jsonreqresp.TicketPurchaseResponse{
		TxID:          t.ticketPurchase.id,
		CustomerName:  t.ticketPurchase.customerName,
		CustomerEmail: t.ticketPurchase.customerEmail,
		PurchaseDate:  t.ticketPurchase.purchaseDate,
		EventID:       t.ticketPurchase.eventID,
		UserID:        t.ticketPurchase.userID,
	}

	return jsonreqresp.TxTicketPurchaseResponse{
		TicketPurchase: ticketPurchaseJson,
		CntTickets:     t.cntTickets,
		ExpiredAt:      t.expiredAt,
	}
}

func (t *TicketPurchaseTx) GetID() uuid.UUID {
	return t.ticketPurchase.GetID()
}

func (t *TicketPurchaseTx) GetExpiredAt() time.Time {
	return t.expiredAt
}

func (t *TicketPurchaseTx) GetTicketPurchase() *TicketPurchase {
	return &t.ticketPurchase
}

func (t *TicketPurchaseTx) GetCntTickets() int {
	return t.cntTickets
}
