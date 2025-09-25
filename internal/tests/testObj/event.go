package testobj

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type EventMother interface {
	EventP(eventID uuid.UUID) *models.Event
	EventCntTicketsP(eventID uuid.UUID, cntTickets int) *models.Event
	EventAdd(employeeID uuid.UUID) *jsonreqresp.EventAdd
	StatCollectionsP() *models.StatCollections
	EventFilterEmpty() *jsonreqresp.EventFilter
}

func NewEventMother() EventMother {
	return &eventMother{}
}

type eventMother struct {
}

func (um *eventMother) EventP(eventID uuid.UUID) *models.Event {
	event, _ := models.NewEvent(
		eventID,
		"test-event-title",
		time.Date(2023, time.October, 1, 15, 30, 0, 0, time.UTC),
		time.Date(2023, time.October, 10, 15, 30, 0, 0, time.UTC),
		"test-adress",
		true,
		uuid.New(),
		100,
		true,
		uuid.UUIDs{uuid.New()}, //, uuid.New(), uuid.New()},
	)
	return &event
}

func (um *eventMother) EventCntTicketsP(eventID uuid.UUID, cntTickets int) *models.Event {
	event, _ := models.NewEvent(
		eventID,
		"test-event-title",
		time.Date(2023, time.October, 1, 15, 30, 0, 0, time.UTC),
		time.Date(2023, time.October, 10, 15, 30, 0, 0, time.UTC),
		"test-adress",
		true,
		uuid.New(),
		cntTickets,
		true,
		uuid.UUIDs{uuid.New(), uuid.New(), uuid.New()},
	)
	return &event
}

func (um *eventMother) EventAdd(employeeID uuid.UUID) *jsonreqresp.EventAdd {
	return &jsonreqresp.EventAdd{
		Title:      "test-event-title-2",
		DateBegin:  time.Date(2023, time.October, 1, 15, 30, 0, 0, time.UTC),
		DateEnd:    time.Date(2023, time.October, 10, 15, 30, 0, 0, time.UTC),
		Address:    "test-adress-2",
		CanVisit:   false,
		EmployeeID: employeeID,
		CntTickets: 200,
		ArtworkIDs: []string{uuid.New().String()},
	}
}

func (um *eventMother) StatCollectionsP() *models.StatCollections {
	st, _ := models.NewStatCollections(
		uuid.New(),
		"test-collection-title",
		10,
	)
	return &st
}

func (um *eventMother) EventFilterEmpty() *jsonreqresp.EventFilter {
	return &jsonreqresp.EventFilter{
		Title:     "",
		DateBegin: time.Time{},
		DateEnd:   time.Time{},
		CanVisit:  "",
		Valid:     "",
	}
}
