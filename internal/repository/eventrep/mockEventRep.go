package eventrep

import (
	"context"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockEventRep реализует EventRep интерфейс для тестирования
type MockEventRep struct {
	mock.Mock
}

func (m *MockEventRep) GetAll(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error) {
	args := m.Called(ctx, filterOps)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRep) GetArtworkIDs(ctx context.Context, eventID uuid.UUID) (uuid.UUIDs, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUIDs), args.Error(1)
}

func (m *MockEventRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRep) GetEventsOfArtworkOnDate(ctx context.Context, artworkID uuid.UUID, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	args := m.Called(ctx, artworkID, dateBeg, dateEnd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRep) CheckEmployeeByID(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockEventRep) Add(ctx context.Context, e *models.Event) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockEventRep) Delete(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventRep) RealDelete(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventRep) Update(ctx context.Context, eventID uuid.UUID, funcUpdate func(*models.Event) (*models.Event, error)) error {
	args := m.Called(ctx, eventID, funcUpdate)
	return args.Error(0)
}

func (m *MockEventRep) AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkIDs uuid.UUIDs) error {
	args := m.Called(ctx, eventID, artworkIDs)
	return args.Error(0)
}

func (m *MockEventRep) DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error {
	args := m.Called(ctx, eventID, artworkID)
	return args.Error(0)
}

func (m *MockEventRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockEventRep) Close() {
	m.Called()
}
