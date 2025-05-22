package auth

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthZ реализует AuthZ интерфейс для тестирования
type MockAuthZ struct {
	mock.Mock
}

func (m *MockAuthZ) Authorize(ctx context.Context, payload token.Payload) context.Context {
	args := m.Called(ctx, payload)
	return args.Get(0).(context.Context)
}

func (m *MockAuthZ) UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	args := m.Called(ctx)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAuthZ) EmployeeIDFromContext(ctx context.Context) (uuid.UUID, error) {
	args := m.Called(ctx)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAuthZ) AdminIDFromContext(ctx context.Context) (uuid.UUID, error) {
	args := m.Called(ctx)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
