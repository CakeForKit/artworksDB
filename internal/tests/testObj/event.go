package testobj

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type EventMother interface {
	EventP(eventID uuid.UUID) *models.Event
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
		time.Now(),
		time.Now().Add(time.Hour*60),
		"test-adress",
		true,
		uuid.New(),
		100,
		true,
		[]uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
	)
	return &event
}
