package eventrep

import (
	"context"
	"errors"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep/mockeventrep"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("the Event was not found in the repository")
)

type EventRep interface {
	GetAll(ctx context.Context) ([]*models.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetByDate(ctx context.Context, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error)
	GetEventOfArtworkOnDate(ctx context.Context, artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error)
	//
	Add(ctx context.Context, aw *models.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error)
}

func NewEventRep() EventRep {
	return &mockeventrep.MockEventRep{}
}
