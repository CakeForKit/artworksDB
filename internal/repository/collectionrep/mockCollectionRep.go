package collectionrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCollectionRep реализует CollectionRep интерфейс для тестирования
type MockCollectionRep struct {
	mock.Mock
}

func (m *MockCollectionRep) GetAllCollections(ctx context.Context) ([]*models.Collection, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Collection), args.Error(1)
}

func (m *MockCollectionRep) GetCollectionByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Collection), args.Error(1)
}

func (m *MockCollectionRep) AddCollection(ctx context.Context, c *models.Collection) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCollectionRep) DeleteCollection(ctx context.Context, idCol uuid.UUID) error {
	args := m.Called(ctx, idCol)
	return args.Error(0)
}

func (m *MockCollectionRep) UpdateCollection(ctx context.Context, idCol uuid.UUID, funcUpdate func(*models.Collection) (*models.Collection, error)) error {
	args := m.Called(ctx, idCol, funcUpdate)
	return args.Error(0)
}
