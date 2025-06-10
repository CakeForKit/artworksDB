package userrep

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("the User was not found in the repository")
	ErrFailedToAddUser    = errors.New("failed to add the User to the repository")
	ErrDuplicateLoginUser = errors.New("a user with this login already exists")
	ErrUpdateUser         = errors.New("failed to update the User in the repository")
)

type UserRep interface {
	GetAll(ctx context.Context) ([]*models.User, error)
	GetAllSubscribed(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	Add(ctx context.Context, e *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error)
	UpdateSubscribeToMailing(ctx context.Context, id uuid.UUID, newSubscribeMail bool) error
	Ping(ctx context.Context) error
	Close()
}

func NewUserRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (UserRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgUserRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHUserRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewUserRep: %w", cnfg.ErrUnknownDB)
	}
	// return &MockAdminRep{}, nil
}
