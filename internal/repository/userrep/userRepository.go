package userrep

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep/mockuserrep"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound    = errors.New("the User was not found in the repository")
	ErrFailedToAddUser = errors.New("failed to add the User to the repository")
	ErrUpdateUser      = errors.New("failed to update the User in the repository")
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

func NewUserRep() UserRep {
	return &mockuserrep.MockUserRep{}
}
