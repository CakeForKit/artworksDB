package mockeventrep

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
)

type MockEventRep struct {
	mock.Mock
}

func (m *MockEventRep) GetAll() []*models.Event {
	args := m.Called()
	return args.Get(0).([]*models.Event)
}

func (m *MockEventRep) GetByID(id uuid.UUID) (*models.Event, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRep) GetByDate(dateBeg time.Time, dateEnd time.Time) []*models.Event {
	args := m.Called(dateBeg, dateEnd)
	return args.Get(0).([]*models.Event)
}

func (m *MockEventRep) GetEventOfArtworkOnDate(artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error) {
	args := m.Called(artwork, dateBeg, dateEnd)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRep) Add(aw *models.Event) error {
	args := m.Called(aw)
	return args.Error(0)
}

func (m *MockEventRep) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventRep) Update(id uuid.UUID,
	funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error) {

	args := m.Called(id, funcUpdate)
	return args.Get(0).(*models.Event), args.Error(1)
}
