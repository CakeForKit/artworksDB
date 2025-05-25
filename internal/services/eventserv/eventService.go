package eventserv

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"github.com/google/uuid"
)

type EventService interface {
	GetAll(ctx context.Context) ([]*models.Event, error)
	GetArtworksFromEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Artwork, error)
	Add(ctx context.Context, eventReq *jsonreqresp.EventAdd) error
	Delete(ctx context.Context, eventID uuid.UUID) error
	Update(ctx context.Context, eventID uuid.UUID, updateFields *jsonreqresp.EventUpdate) error
	AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkIDs uuid.UUIDs) error
	DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error
}

var (
	ErrArtworkBusy = errors.New("artowrk can't participate in event")
)

type eventService struct {
	eventRep   eventrep.EventRep
	artworkRep artworkrep.ArtworkRep
}

func NewEventService(eventRep eventrep.EventRep, artworkRep artworkrep.ArtworkRep) EventService {
	return &eventService{
		eventRep:   eventRep,
		artworkRep: artworkRep,
	}
}

func (e *eventService) GetAll(ctx context.Context) ([]*models.Event, error) {
	return e.eventRep.GetAll(ctx, &jsonreqresp.EventFilter{})
}

func (e *eventService) GetArtworksFromEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Artwork, error) {
	artworkIDs, err := e.eventRep.GetArtworkIDs(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("eventService.GetArtworkFromEvent: %w", err)
	}
	artworks := make([]*models.Artwork, len(artworkIDs))
	for i, aID := range artworkIDs {
		art, err := e.artworkRep.GetByID(ctx, aID)
		if err != nil {
			return nil, fmt.Errorf("eventService.GetArtworkFromEvent: %w", err)
		}
		artworks[i] = art
	}
	return artworks, nil
}

func (e *eventService) Add(ctx context.Context, eventReq *jsonreqresp.EventAdd) error {
	employeeExist, err := e.eventRep.CheckEmployeeByID(ctx, eventReq.EmployeeID)
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
		eventReq.EmployeeID,
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
