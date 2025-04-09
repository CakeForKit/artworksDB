package userrep

import (
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
	GetAll() []*models.User
	GetAllSubscribed() []*models.User
	GetByID(id uuid.UUID) (*models.User, error)
	GetByLogin(login string) (*models.User, error)
	Add(e *models.User) error
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error)
}

func NewUserRep() UserRep {
	return &mockuserrep.MockUserRep{}
}
