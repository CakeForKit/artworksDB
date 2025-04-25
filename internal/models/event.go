package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	id         uuid.UUID
	title      string
	dateBegin  time.Time
	dateEnd    time.Time
	address    string
	access     bool
	employeeID uuid.UUID
	cntTickets int
}

var (
	ErrEventEmptyTitle      = errors.New("empty title")
	ErrEventTitleTooLong    = errors.New("title exceeds maximum length (255 chars)")
	ErrEventInvalidDates    = errors.New("end date must be after start date")
	ErrEventEmptyAddress    = errors.New("empty address")
	ErrEventAddressTooLong  = errors.New("address exceeds maximum length (255 chars)")
	ErrEventInvalidAccess   = errors.New("invalid access value")
	ErrEventInvalidEmployee = errors.New("invalid employee ID")
	ErrEventNegativeTickets = errors.New("ticket count cannot be negative")
)

func NewEvent(
	id uuid.UUID,
	title string,
	dateBegin time.Time,
	dateEnd time.Time,
	address string,
	access bool,
	employeeID uuid.UUID,
	cntTickets int,
) (Event, error) {
	event := Event{
		id:         id,
		title:      strings.TrimSpace(title),
		dateBegin:  dateBegin,
		dateEnd:    dateEnd,
		address:    strings.TrimSpace(address),
		access:     access,
		employeeID: employeeID,
		cntTickets: cntTickets,
	}

	if err := event.validate(); err != nil {
		return Event{}, err
	}

	return event, nil
}

func (e *Event) validate() error {
	switch {
	case e.title == "":
		return ErrEventEmptyTitle
	case len(e.title) > 255:
		return ErrEventTitleTooLong
	case e.dateEnd.Before(e.dateBegin):
		return ErrEventInvalidDates
	case e.address == "":
		return ErrEventEmptyAddress
	case len(e.address) > 255:
		return ErrEventAddressTooLong
	case e.employeeID == uuid.Nil:
		return ErrEventInvalidEmployee
	case e.cntTickets < 0:
		return ErrEventNegativeTickets
	}
	return nil
}

func (e *Event) TextAbout() string {
	return fmt.Sprintf(
		"Событие: %s\nДата начала: %s\nДата окончания: %s\nАдрес: %s\nДоступ: %v\n"+
			"Ответственный сотрудник: %s\nДоступно билетов: %d",
		e.GetTitle(),
		e.GetDateBegin().Format(time.RFC3339),
		e.GetDateEnd().Format(time.RFC3339),
		e.GetAddress(),
		e.GetAccess(),
		e.GetEmployeeID(),
		e.GetTicketCount(),
	)
}

func (e *Event) GetID() uuid.UUID {
	return e.id
}

func (e *Event) GetTitle() string {
	return e.title
}

func (e *Event) GetDateBegin() time.Time {
	return e.dateBegin
}

func (e *Event) GetDateEnd() time.Time {
	return e.dateEnd
}

func (e *Event) GetAddress() string {
	return e.address
}

func (e *Event) GetAccess() bool {
	return e.access
}

func (e *Event) GetEmployeeID() uuid.UUID {
	return e.employeeID
}

func (e *Event) GetTicketCount() int {
	return e.cntTickets
}
