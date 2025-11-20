package adminserv

import (
	"context"
	"fmt"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/authz"
	"github.com/google/uuid"
)

type AdminService interface {
	GetAllEmployees(ctx context.Context) ([]*models.Employee, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	ChangeEmployeeRights(ctx context.Context, employeeID uuid.UUID, valid bool) error
}

func NewAdminService(empRep employeerep.EmployeeRep, userRep userrep.UserRep, authZ authz.AuthZ) AdminService {
	return &adminService{
		employeeRep: empRep,
		userRep:     userRep,
		authZ:       authZ,
	}
}

type adminService struct {
	employeeRep employeerep.EmployeeRep
	userRep     userrep.UserRep
	authZ       authz.AuthZ
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
		return fmt.Errorf("adminService.ChangeRights: %w", err)
	}
	return nil
}
