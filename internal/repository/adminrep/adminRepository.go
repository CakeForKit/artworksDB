package adminrep

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrAdminNotFound = errors.New("the Admin was not found in the repository")
	// ErrFailedToAddAdmin  = errors.New("failed to add the Admin to the repository")
	ErrDuplicateLoginAdm = errors.New("an admin with this login already exists")
	ErrUpdateAdmin       = errors.New("failed to update the Admin in the repository")
)

type AdminRep interface {
	GetAll(ctx context.Context) ([]*models.Admin, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Admin, error)
	GetByLogin(ctx context.Context, login string) (*models.Admin, error)
	Add(ctx context.Context, e *models.Admin) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Admin) (*models.Admin, error)) error
	Ping(ctx context.Context) error
	Close()
}

func NewAdminRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (AdminRep, error) {
	return NewPgAdminRep(ctx, pgCreds, dbConf)
	// return &MockAdminRep{}, nil
}
