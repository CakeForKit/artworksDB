package eventrep

import (
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
	GetAll() []*models.Event
	GetByID(uuid.UUID) (*models.Event, error)
	GetByDate(dateBeg time.Time, dateEnd time.Time) []*models.Event
	GetEventOfArtworkOnDate(artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error)
	//
	Add(aw *models.Event) error
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error)
}

func NewEventRep() EventRep {
	return &mockeventrep.MockEventRep{}
}
