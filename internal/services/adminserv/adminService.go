package adminserv

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/google/uuid"
)

type AdminService interface {
	GetAllEmployees(ctx context.Context) ([]*models.Employee, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	ChangeEmployeeRights(ctx context.Context, employeeID uuid.UUID, valid bool) error
}

func NewAdminService(empRep employeerep.EmployeeRep, userRep userrep.UserRep, authZ auth.AuthZ) AdminService {
	return &adminService{
		employeeRep: empRep,
		userRep:     userRep,
		authZ:       authZ,
	}
}

type adminService struct {
	employeeRep employeerep.EmployeeRep
	userRep     userrep.UserRep
	authZ       auth.AuthZ
}

func (e *adminService) GetAllEmployees(ctx context.Context) ([]*models.Employee, error) {
	_, err := e.authZ.AdminIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("adminService.GetAllEmployees: %w", err)
	}
	return e.employeeRep.GetAll(ctx)
}

func (e *adminService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	_, err := e.authZ.AdminIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("adminService.GetAllEmployees: %w", err)
	}
	return e.userRep.GetAll(ctx)
}

func (e *adminService) ChangeEmployeeRights(ctx context.Context, employeeID uuid.UUID, valid bool) error {
	_, err := e.authZ.AdminIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("adminService.GetAllEmployees: %w", err)
	}

	funcUpdate := func(empl *models.Employee) (*models.Employee, error) {
		empl.SetValid(valid)
		return empl, nil
	}
	_, err = e.employeeRep.Update(ctx, employeeID, funcUpdate)
	if err != nil {
		return fmt.Errorf("mployeeService.ChangeRights: %w", err)
	}
	return nil
}
