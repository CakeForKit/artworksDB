package eventrep

import (
	"context"
	"errors"
	"fmt"
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
	ErrUpdateEvent          = errors.New("err update Event params")
	// ErrUpdateNoEmployee     = errors.New("failed to update the Events, no employeee")
)

type EventRep interface {
	GetAll(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error)
	GetArtworkIDs(ctx context.Context, eventID uuid.UUID) (uuid.UUIDs, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetEventsOfArtworkOnDate(ctx context.Context, artworkID uuid.UUID, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error)
	GetCollectionsStat(ctx context.Context, eventID uuid.UUID) ([]*models.StatCollections, error)
	CheckEmployeeByID(ctx context.Context, id uuid.UUID) (bool, error)
	//
	Add(ctx context.Context, e *models.Event) error
	Delete(ctx context.Context, eventID uuid.UUID) error
	RealDelete(ctx context.Context, eventID uuid.UUID) error
	Update(ctx context.Context, eventID uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) error
	AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUIDs) error
	DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error
	Ping(ctx context.Context) error
	Close()
}

func NewEventRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (EventRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgEventRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHEventRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewEventRep: %w", cnfg.ErrUnknownDB)
	}
	// return &MockAdminRep{}, nil
}
