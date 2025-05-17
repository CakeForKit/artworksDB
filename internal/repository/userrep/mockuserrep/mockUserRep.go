package mockuserrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
)

type MockUserRep struct {
	mock.Mock
}

func (m *MockUserRep) GetAll(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRep) GetAllSubscribed(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRep) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRep) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRep) Add(ctx context.Context, e *models.User) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockUserRep) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRep) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error) {
	args := m.Called(ctx, id, funcUpdate)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRep) UpdateSubscribeToMailing(ctx context.Context, id uuid.UUID, newSubscribeMail bool) error {
	args := m.Called(ctx, id, newSubscribeMail)
	return args.Error(0)
}

func (m *MockUserRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserRep) Close() {

}
