package eventserv

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"github.com/google/uuid"
)

type EventService interface {
	GetAllEvents(ctx context.Context) ([]*models.Event, error)
	Add(ctx context.Context, e *models.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error)
}

type eventService struct {
	eventRep eventrep.EventRep
}

func NewEventService(artRep eventrep.EventRep) EventService {
	return &eventService{
		eventRep: artRep,
	}
}

func (e *eventService) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	return e.eventRep.GetAll(ctx)
}

func (e *eventService) Add(ctx context.Context, aw *models.Event) error {
	return e.eventRep.Add(ctx, aw)
}

func (e *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return e.eventRep.Delete(ctx, id)
}

func (e *eventService) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error) {
	return e.eventRep.Update(ctx, id, funcUpdate)
}
