package models

import (
	"errors"
	"strings"
	"time"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type TicketPurchase struct {
	id            uuid.UUID
	customerName  string
	customerEmail string
	purchaseDate  time.Time
	eventID       uuid.UUID
	userID        uuid.UUID
}

type jsonTicketPurchase struct {
	ID            uuid.UUID `json:"id"`
	CustomerName  string    `json:"customerName"`
	CustomerEmail string    `json:"customerEmail"`
	PurchaseDate  time.Time `json:"purchaseDate"`
	EventID       uuid.UUID `json:"eventId"`
	UserID        uuid.UUID `json:"userId"`
}

var (
	ErrTicketPurchaseEmptyName    = errors.New("empty customer name")
	ErrTicketPurchaseNameTooLong  = errors.New("customer name exceeds maximum length (100 chars)")
	ErrTicketPurchaseEmailTooLong = errors.New("customer email exceeds maximum length (100 chars)")
	ErrTicketPurchaseInvalidEmail = errors.New("invalid customer email")
	ErrTicketPurchaseEmptyEventID = errors.New("empty event ID")
	ErrTicketPurchaseInvalidDate  = errors.New("invalid purchase date")
)

func NewTicketPurchase(
	id uuid.UUID,
	customerName string,
	customerEmail string,
	purchaseDate time.Time,
	eventID uuid.UUID,
	userID uuid.UUID,
) (TicketPurchase, error) {
	tp := TicketPurchase{
		id:            id,
		customerName:  strings.TrimSpace(customerName),
		customerEmail: strings.TrimSpace(customerEmail),
		purchaseDate:  purchaseDate,
		eventID:       eventID,
		userID:        userID,
	}

	if err := tp.validate(); err != nil {
		return TicketPurchase{}, err
	}

	return tp, nil
}

func (tp *TicketPurchase) validate() error {
	switch {
	case tp.customerName == "":
		return ErrTicketPurchaseEmptyName
	case len(tp.customerName) > 100:
		return ErrTicketPurchaseNameTooLong
	case len(tp.customerEmail) > 100:
		return ErrTicketPurchaseEmailTooLong
	case !isValidEmail(tp.customerEmail):
		return ErrTicketPurchaseInvalidEmail
	case tp.eventID == uuid.Nil:
		return ErrTicketPurchaseEmptyEventID
	case tp.purchaseDate.IsZero():
		return ErrTicketPurchaseInvalidDate
	}
	return nil
}

func (tp *TicketPurchase) GetID() uuid.UUID {
	return tp.id
}

func (tp *TicketPurchase) GetCustomerName() string {
	return tp.customerName
}

func (tp *TicketPurchase) GetCustomerEmail() string {
	return tp.customerEmail
}

func (tp *TicketPurchase) GetPurchaseDate() time.Time {
	return tp.purchaseDate
}

func (tp *TicketPurchase) GetEventID() uuid.UUID {
	return tp.eventID
}

func (tp *TicketPurchase) GetUserID() uuid.UUID {
	return tp.userID
}

func (t *TicketPurchase) ToTicketPurchaseResponse() jsonreqresp.TicketPurchaseResponse {
	return jsonreqresp.TicketPurchaseResponse{
		ID:            t.id,
		CustomerName:  t.customerName,
		CustomerEmail: t.customerEmail,
		PurchaseDate:  t.purchaseDate,
		EventID:       t.eventID,
		UserID:        t.userID,
	}
}

// // SetCustomerName устанавливает имя покупателя
// func (tp *TicketPurchase) SetCustomerName(name string) error {
// 	name = strings.TrimSpace(name)
// 	if name == "" {
// 		return ErrTicketPurchaseEmptyName
// 	}
// 	if len(name) > 100 {
// 		return ErrTicketPurchaseNameTooLong
// 	}
// 	tp.customerName = name
// 	return nil
// }

// // SetCustomerEmail устанавливает email покупателя
// func (tp *TicketPurchase) SetCustomerEmail(email string) error {
// 	email = strings.TrimSpace(email)
// 	if !isValidEmail(email) {
// 		return ErrTicketPurchaseInvalidEmail
// 	}
// 	tp.customerEmail = email
// 	return nil
// }

// // SetPurchaseDate устанавливает дату покупки
// func (tp *TicketPurchase) SetPurchaseDate(date time.Time) error {
// 	if date.IsZero() {
// 		return ErrTicketPurchaseInvalidDate
// 	}
// 	tp.purchaseDate = date
// 	return nil
// }

// // SetEventID устанавливает ID события
// func (tp *TicketPurchase) SetEventID(eventID uuid.UUID) error {
// 	if eventID == uuid.Nil {
// 		return ErrTicketPurchaseEmptyEventID
// 	}
// 	tp.eventID = eventID
// 	return nil
// }

// // SetUserID устанавливает ID пользователя
// func (tp *TicketPurchase) SetUserID(userID uuid.UUID) error {
// 	if userID == uuid.Nil {
// 		return ErrTicketPurchaseEmptyUserID
// 	}
// 	tp.userID = userID
// 	return nil
// }
