package employeeserv

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"github.com/google/uuid"
)

type EmployeeService interface {
	GetAllEmployees() []*models.Employee
	// Add(*models.Employee) error // нехешированный пароль (здесь он и хешируется)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error)
}

type employeeService struct {
	employeeRep employeerep.EmployeeRep
}

func NewEmployeeService(empRep employeerep.EmployeeRep) EmployeeService {
	return &employeeService{
		employeeRep: empRep,
	}
}

func (e *employeeService) GetAllEmployees() []*models.Employee {
	return e.employeeRep.GetAll()
}

// func (e *employeeService) Add(emp *models.Employee) error {
// 	return e.employeeRep.Add(emp)
// }

func (e *employeeService) Delete(id uuid.UUID) error {
	return e.employeeRep.Delete(id)
}

func (e *employeeService) Update(id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error) {
	return e.employeeRep.Update(id, funcUpdate)
}
