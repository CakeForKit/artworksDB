package eventserv

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"github.com/google/uuid"
)

type EventService interface {
	GetAllEvents() []*models.Event
	Add(*models.Event) error
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error)
}

type eventService struct {
	eventRep eventrep.EventRep
}

func NewEventService(artRep eventrep.EventRep) EventService {
	return &eventService{
		eventRep: artRep,
	}
}

func (e *eventService) GetAllEvents() []*models.Event {
	return e.eventRep.GetAll()
}

func (e *eventService) Add(aw *models.Event) error {
	return e.eventRep.Add(aw)
}

func (e *eventService) Delete(id uuid.UUID) error {
	return e.eventRep.Delete(id)
}

func (e *eventService) Update(id uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error) {
	return e.eventRep.Update(id, funcUpdate)
}
