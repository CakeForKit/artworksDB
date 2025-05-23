package eventrep

import (
	"context"
	"errors"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound        = errors.New("the Event was not found in the repository")
	ErrEventArtowrkNotFound = errors.New("the Event_artwork was not found in the repository")
	ErrAddNoEmployee        = errors.New("failed to add the Event, no employeee")
	ErrUpdateNoEmployee     = errors.New("failed to update the Events, no employeee")
)

type EventRep interface {
	GetAll(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error)
	GetArtworkIDs(ctx context.Context, eventID uuid.UUID) (uuid.UUIDs, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetEventsOfArtworkOnDate(ctx context.Context, artworkID uuid.UUID, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error)
	CheckEmployeeByID(ctx context.Context, id uuid.UUID) (bool, error)
	//
	Add(ctx context.Context, e *models.Event) error
	Delete(ctx context.Context, eventID uuid.UUID) error
	Update(ctx context.Context, eventID uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) error
	AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUIDs) error
	DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error
	Ping(ctx context.Context) error
	Close()
}

func NewEventRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (EventRep, error) {
	rep, err := NewPgEventRep(ctx, pgCreds, dbConf)
	return rep, err
}
