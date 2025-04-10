package mockemployeerep

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
)

type MockEmployeeRep struct {
	mock.Mock
}

func (m *MockEmployeeRep) GetAll() []*models.Employee {
	args := m.Called()
	return args.Get(0).([]*models.Employee)
}

func (m *MockEmployeeRep) GetByLogin(login string) (*models.Employee, error) {
	args := m.Called(login)
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRep) Add(e *models.Employee) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockEmployeeRep) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeRep) Update(id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error) {
	args := m.Called(id, funcUpdate)
	return args.Get(0).(*models.Employee), args.Error(1)
}
