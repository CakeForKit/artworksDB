package employeeserv

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"github.com/google/uuid"
)

type EmployeeService interface {
	GetAllEmployees(ctx context.Context) ([]*models.Employee, error)
	ChangeRights(ctx context.Context, employeeID uuid.UUID, valid bool) error
	// Delete(ctx context.Context, id uuid.UUID) error
	// Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error)
}

func NewEmployeeService(empRep employeerep.EmployeeRep) EmployeeService {
	return &employeeService{
		employeeRep: empRep,
	}
}

type employeeService struct {
	employeeRep employeerep.EmployeeRep
}

func (e *employeeService) GetAllEmployees(ctx context.Context) ([]*models.Employee, error) {
	return e.employeeRep.GetAll(ctx)
}

func (e *employeeService) ChangeRights(ctx context.Context, employeeID uuid.UUID, valid bool) error {
	funcUpdate := func(empl *models.Employee) (*models.Employee, error) {
		empl.SetValid(valid)
		return empl, nil
	}
	_, err := e.employeeRep.Update(ctx, employeeID, funcUpdate)
	if err != nil {
		return fmt.Errorf("mployeeService.ChangeRights: %w", err)
	}
	return nil
}

// func (e *employeeService) Delete(ctx context.Context, id uuid.UUID) error {
// 	return e.employeeRep.Delete(ctx, id)
// }

// func (e *employeeService) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error) {
// 	return e.employeeRep.Update(ctx, id, funcUpdate)
// }
