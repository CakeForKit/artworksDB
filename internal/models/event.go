package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type Event struct {
	id         uuid.UUID
	title      string
	dateBegin  time.Time
	dateEnd    time.Time
	address    string
	canVisit   bool
	employeeID uuid.UUID
	cntTickets int
	artworkIDs uuid.UUIDs
	valid      bool
}

var (
	ErrValidateEvent        = errors.New("invalid model Event")
	ErrAddArtwork           = errors.New("the artowrk is already participating in the event")
	ErrDeleteArtwork        = errors.New("the artwork was not in event")
	ErrEventEmptyTitle      = errors.New("empty title")
	ErrEventTitleTooLong    = errors.New("title exceeds maximum length (255 chars)")
	ErrEventInvalidDates    = errors.New("end date must be after start date")
	ErrEventEmptyAddress    = errors.New("empty address")
	ErrEventAddressTooLong  = errors.New("address exceeds maximum length (255 chars)")
	ErrEventInvalidAccess   = errors.New("invalid canVisit value")
	ErrEventInvalidEmployee = errors.New("invalid employee ID")
	ErrEventNegativeTickets = errors.New("ticket count cannot be negative")
	ErrDuplicateArtwokIDs   = errors.New("duplicate artwork ids")
)

func NewEvent(
	id uuid.UUID,
	title string,
	dateBegin time.Time,
	dateEnd time.Time,
	address string,
	canVisit bool,
	employeeID uuid.UUID,
	cntTickets int,
	valid bool,
	artworkIDs uuid.UUIDs,
) (Event, error) {
	event := Event{
		id:         id,
		title:      strings.TrimSpace(title),
		dateBegin:  dateBegin,
		dateEnd:    dateEnd,
		address:    strings.TrimSpace(address),
		canVisit:   canVisit,
		employeeID: employeeID,
		cntTickets: cntTickets,
		artworkIDs: artworkIDs,
		valid:      valid,
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
	case HasDuplicateUUIDs(e.artworkIDs):
		return ErrDuplicateArtwokIDs
	}
	return nil
}

func HasDuplicateUUIDs(ids uuid.UUIDs) bool {
	seen := make(map[uuid.UUID]struct{})
	for _, id := range ids {
		if _, exists := seen[id]; exists {
			return true
		}
		seen[id] = struct{}{}
	}
	return false
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
	return e.canVisit
}

func (e *Event) GetEmployeeID() uuid.UUID {
	return e.employeeID
}

func (e *Event) GetTicketCount() int {
	return e.cntTickets
}

func (e *Event) GetArtworkIDs() uuid.UUIDs {
	return e.artworkIDs
}

func (e *Event) IsValid() bool {
	return e.valid
}

func (e *Event) AddArtworks(idArts uuid.UUIDs) error {
	for _, oldID := range e.artworkIDs {
		for _, newID := range idArts {
			if oldID == newID {
				return fmt.Errorf("Event.AddArtworks %w %w", ErrValidateEvent, ErrAddArtwork)
			}
		}

	}
	e.artworkIDs = append(e.artworkIDs, idArts...)
	return nil
}

func (e *Event) DeleteArtwork(idArt uuid.UUID) error {
	for i, v := range e.artworkIDs {
		if v == idArt {
			e.artworkIDs = append(e.artworkIDs[:i], e.artworkIDs[i+1:]...)
			return nil
		}
	}
	return ErrDeleteArtwork
}

func (e *Event) Update(updateReq *jsonreqresp.EventUpdate) error {
	copyE := *e
	copyE.title = updateReq.Title
	copyE.dateBegin = updateReq.DateBegin
	copyE.dateEnd = updateReq.DateEnd
	copyE.address = updateReq.Address
	copyE.canVisit = updateReq.CanVisit
	copyE.cntTickets = updateReq.CntTickets
	copyE.valid = updateReq.Valid

	if err := copyE.validate(); err != nil {
		return err
	}
	*e = copyE
	return nil
}

func (e *Event) ToEventResponse() jsonreqresp.EventResponse {
	return jsonreqresp.EventResponse{
		ID:         e.id.String(),
		Title:      e.title,
		DateBegin:  e.dateBegin,
		DateEnd:    e.dateEnd,
		Address:    e.address,
		CanVisit:   e.canVisit,
		EmployeeID: e.employeeID.String(),
		CntTickets: e.cntTickets,
		Valid:      e.valid,
		ArtworkIDs: e.GetArtworkIDs().Strings(),
	}
}
