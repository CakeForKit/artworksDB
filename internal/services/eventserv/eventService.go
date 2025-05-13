package eventserv

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"github.com/google/uuid"
)

type EventService interface {
	GetAll(ctx context.Context) ([]*models.Event, error)
	Add(ctx context.Context, eventReq *jsonreqresp.AddEventRequest) error
	Delete(ctx context.Context, eventID uuid.UUID) error
	Update(ctx context.Context, eventID uuid.UUID, updateFields *jsonreqresp.EventUpdate) error
	AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkIDs uuid.UUIDs) error
	DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error
}

var (
	ErrArtworkBusy = errors.New("artowrk can't participate in event")
)

type eventService struct {
	eventRep eventrep.EventRep
}

func NewEventService(eventRep eventrep.EventRep) EventService {
	return &eventService{
		eventRep: eventRep,
	}
}

func (e *eventService) GetAll(ctx context.Context) ([]*models.Event, error) {
	return e.eventRep.GetAll(ctx)
}

func (e *eventService) Add(ctx context.Context, eventReq *jsonreqresp.AddEventRequest) error {
	employeeID := uuid.MustParse(eventReq.EmployeeID)
	employeeExist, err := e.eventRep.CheckEmployeeByID(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("eventService.Add check employee: %v", err)
	} else if !employeeExist {
		return fmt.Errorf("eventService.Add check employee: %w", eventrep.ErrAddNoEmployee)
	}

	var artworkIDs uuid.UUIDs
	for _, v := range eventReq.ArtworkIDs {
		artworkIDs = append(artworkIDs, uuid.MustParse(v))
	}
	event, err := models.NewEvent(
		uuid.New(),
		eventReq.Title,
		eventReq.DateBegin,
		eventReq.DateEnd,
		eventReq.Address,
		eventReq.CanVisit,
		employeeID,
		eventReq.CntTickets,
		true,
		artworkIDs,
	)
	if err != nil {
		return fmt.Errorf("eventService.Add %w: %v", models.ErrValidateEvent, err)
	}

	for _, id := range artworkIDs {
		_, err := e.eventRep.GetEventsOfArtworkOnDate(ctx, id, eventReq.DateBegin, eventReq.DateEnd)
		if err == nil {
			return fmt.Errorf("eventService.Add: %w", ErrArtworkBusy)
		} else if err != eventrep.ErrEventNotFound {
			return fmt.Errorf("eventService.Add: %w", err)
		}
	}

	err = e.eventRep.Add(ctx, &event)
	if err != nil {
		return fmt.Errorf("eventService.Add: %v", err)
	}
	err = e.eventRep.AddArtworksToEvent(ctx, event.GetID(), artworkIDs)
	if err != nil {
		return fmt.Errorf("eventService.Add: %v", err)
	}
	return nil
}

func (e *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return e.eventRep.Delete(ctx, id)
}

func (e *eventService) Update(ctx context.Context, eventID uuid.UUID, updateFields *jsonreqresp.EventUpdate) error {
	_, err := e.eventRep.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("eventService.Update: %v", err)
	}

	return e.eventRep.Update(
		ctx,
		eventID,
		func(event *models.Event) (*models.Event, error) {
			err := event.Update(updateFields)
			return event, err
		})
}

func (e *eventService) AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkIDs uuid.UUIDs) error {
	if models.HasDuplicateUUIDs(artworkIDs) {
		return fmt.Errorf("PgEventRep.AddArtworkToEvent: %v", models.ErrDuplicateArtwokIDs)
	}
	_, err := e.eventRep.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("PgEventRep.AddArtworkToEvent: %v", err)
	}
	oldArtworksIDs, err := e.eventRep.GetArtworkIDs(ctx, eventID)
	if err != nil {
		return fmt.Errorf("PgEventRep.AddArtworkToEvent: %v", err)
	}
	for _, oldID := range oldArtworksIDs {
		for _, newID := range artworkIDs {
			if oldID == newID {
				return fmt.Errorf("gEventRep.AddArtworkToEvent: %w", models.ErrAddArtwork)
			}
		}

	}
	return e.eventRep.AddArtworksToEvent(ctx, eventID, artworkIDs)
}

func (e *eventService) DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error {
	return e.eventRep.DeleteArtworkFromEvent(ctx, eventID, artworkID)
}
