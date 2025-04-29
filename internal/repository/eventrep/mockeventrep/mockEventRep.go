package mockeventrep

import (
	"context"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
)

type MockEventRep struct {
	mock.Mock
}

func (m *MockEventRep) GetAll(ctx context.Context) ([]*models.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRep) GetByDate(ctx context.Context, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	args := m.Called(ctx, dateBeg, dateEnd)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRep) GetEventOfArtworkOnDate(ctx context.Context, artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error) {
	args := m.Called(ctx, artwork, dateBeg, dateEnd)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRep) Add(ctx context.Context, aw *models.Event) error {
	args := m.Called(ctx, aw)
	return args.Error(0)
}

func (m *MockEventRep) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRep) Update(ctx context.Context, id uuid.UUID,
	funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error) {

	args := m.Called(ctx, id, funcUpdate)
	return args.Get(0).(*models.Event), args.Error(1)
}
