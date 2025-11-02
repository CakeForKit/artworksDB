package collectionrep

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCollectionRep реализует CollectionRep интерфейс для тестирования
type MockCollectionRep struct {
	mock.Mock
}

func (m *MockCollectionRep) GetAll(ctx context.Context) ([]*models.Collection, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Collection), args.Error(1)
}

func (m *MockCollectionRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Collection), args.Error(1)
}

func (m *MockCollectionRep) Add(ctx context.Context, c *models.Collection) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCollectionRep) Delete(ctx context.Context, idCol uuid.UUID) error {
	args := m.Called(ctx, idCol)
	return args.Error(0)
}

func (m *MockCollectionRep) Update(ctx context.Context, idCol uuid.UUID, funcUpdate func(*models.Collection) (*models.Collection, error)) error {
	args := m.Called(ctx, idCol, funcUpdate)
	return args.Error(0)
}

func (m *MockCollectionRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCollectionRep) Close() {
	m.Called()
}
