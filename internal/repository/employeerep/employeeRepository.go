package employeerep

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrEmployeeNotFound    = errors.New("the Employee was not found in the repository")
	ErrDuplicateLoginEmp   = errors.New("an employee with this login already exists")
	ErrFailedToAddEmployee = errors.New("failed to add the Employee to the repository")
	ErrUpdateEmployee      = errors.New("failed to update the Employee in the repository")
)

type EmployeeRep interface {
	GetAll(ctx context.Context) ([]*models.Employee, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetByLogin(ctx context.Context, login string) (*models.Employee, error)
	Add(ctx context.Context, e *models.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error)
	Ping(ctx context.Context) error
	Close()
}

func NewEmployeeRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (EmployeeRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgEmployeeRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHEmployeeRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewEmployeeRep: %w", cnfg.ErrUnknownDB)
	}
	// return &MockAdminRep{}, nil
}
