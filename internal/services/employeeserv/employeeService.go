package employeeserv

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"github.com/google/uuid"
)

type EmployeeService interface {
	GetAllEmployees(ctx context.Context) []*models.Employee
	// Add(*models.Employee) error // нехешированный пароль (здесь он и хешируется)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error)
}

type employeeService struct {
	employeeRep employeerep.EmployeeRep
}

func NewEmployeeService(empRep employeerep.EmployeeRep) EmployeeService {
	return &employeeService{
		employeeRep: empRep,
	}
}

func (e *employeeService) GetAllEmployees(ctx context.Context) []*models.Employee {
	return e.employeeRep.GetAll(ctx)
}

// func (e *employeeService) Add(emp *models.Employee) error {
// 	return e.employeeRep.Add(emp)
// }

func (e *employeeService) Delete(ctx context.Context, id uuid.UUID) error {
	return e.employeeRep.Delete(ctx, id)
}

func (e *employeeService) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error) {
	return e.employeeRep.Update(ctx, id, funcUpdate)
}
