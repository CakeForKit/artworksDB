package adminrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAdminRep реализует AdminRep интерфейс для тестирования
type MockAdminRep struct {
	mock.Mock
}

func (m *MockAdminRep) GetAll(ctx context.Context) ([]*models.Admin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Admin), args.Error(1)
}

func (m *MockAdminRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Admin, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAdminRep) GetByLogin(ctx context.Context, login string) (*models.Admin, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAdminRep) Add(ctx context.Context, a *models.Admin) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAdminRep) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAdminRep) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Admin) (*models.Admin, error)) error {
	args := m.Called(ctx, id, funcUpdate)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func (m *MockAdminRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAdminRep) Close() {
	m.Called()
}
