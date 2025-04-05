package mockuserrep

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
)

type MockUserRep struct {
	mock.Mock
}

func (m *MockUserRep) GetAll() []*models.User {
	args := m.Called()
	return args.Get(0).([]*models.User)
}

func (m *MockUserRep) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRep) GetByLogin(login string) (*models.User, error) {
	args := m.Called(login)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRep) Add(e *models.User) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockUserRep) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRep) Update(id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error) {
	args := m.Called(id, funcUpdate)
	return args.Get(0).(*models.User), args.Error(1)
}
