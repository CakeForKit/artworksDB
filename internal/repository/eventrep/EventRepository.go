package eventrep

import (
	"context"
	"errors"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound    = errors.New("the Event was not found in the repository")
	ErrAddNoEmployee    = errors.New("failed to add the Event, no employeee")
	ErrUpdateNoEmployee = errors.New("failed to update the Events, no employeee")
)

type EventRep interface {
	GetAll(ctx context.Context) ([]*models.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetByDate(ctx context.Context, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error)
	GetEventsOfArtworkOnDate(ctx context.Context, artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error)
	//
	Add(ctx context.Context, aw *models.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error)
	Ping(ctx context.Context) error
	Close()
}

func NewEventRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (EventRep, error) {
	rep, err := NewPgEventRep(ctx, pgCreds, dbConf)
	return rep, err
}
