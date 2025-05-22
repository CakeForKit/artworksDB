package jsonreqresp

import (
	"time"

	"github.com/google/uuid"
)

type BuyTicketRequest struct {
	EventID       string `json:"eventID" binding:"required,uuid" example:"b10f841d-ba75-48df-a9cf-c86fc9bd3041"`
	CntTickets    int    `json:"cntTickets" binding:"required,min=0" example:"1"`
	CustomerName  string `json:"customerName,omitempty" binding:"omitempty,max=100" example:"myname"`
	CustomerEmail string `json:"CustomerEmail,omitempty" binding:"omitempty,max=100" example:"myname@test.ru"`
}

type TxTicketPurchaseResponse struct {
	TicketPurchase TicketPurchaseResponse
	CntTickets     int       `json:"cntTickets"`
	ExpiredAt      time.Time `json:"expiredAt"`
}

type TicketPurchaseResponse struct {
	TxID          uuid.UUID `json:"id"`
	CustomerName  string    `json:"customerName"`
	CustomerEmail string    `json:"customerEmail"`
	PurchaseDate  time.Time `json:"purchaseDate"`
	EventID       uuid.UUID `json:"eventId"`
	UserID        uuid.UUID `json:"userId"`
}

type ConfirmCancelTxRequest struct {
	TxID string `json:"txID" binding:"required,uuid" example:"b10f841d-ba75-48df-a9cf-c86fc9bd3041"`
}
