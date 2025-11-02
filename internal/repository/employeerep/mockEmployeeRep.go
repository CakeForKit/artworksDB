package employeerep

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockEmployeeRep struct {
	mock.Mock
}

func (m *MockEmployeeRep) GetAll(ctx context.Context) ([]*models.Employee, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Employee), args.Error(1)
}

func (m *MockEmployeeRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRep) GetByLogin(ctx context.Context, login string) (*models.Employee, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRep) Add(ctx context.Context, e *models.Employee) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockEmployeeRep) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmployeeRep) Update_(ctx context.Context, updatedEmployee *models.Employee) (*models.Employee, error) {
	args := m.Called(ctx, updatedEmployee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRep) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error) {
	args := m.Called(ctx, id, funcUpdate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockEmployeeRep) Close() {
}
