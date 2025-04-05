package employeerep

import (
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep/mockemployeerep"
	"github.com/google/uuid"
)

var (
	ErrEmployeeNotFound    = errors.New("the Employee was not found in the repository")
	ErrFailedToAddEmployee = errors.New("failed to add the Employee to the repository")
	ErrUpdateEmployee      = errors.New("failed to update the Employee in the repository")
)

type EmployeeRep interface {
	GetAll() []*models.Employee
	GetByLogin(login string) (*models.Employee, error)
	Add(e *models.Employee) error
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error)
}

func NewEmployeeRep() EmployeeRep {
	return &mockemployeerep.MockEmployeeRep{}
}
