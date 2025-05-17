package userservice

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	// GetByLogin(login string) (*models.User, error)
	// Add(*models.User) error // нехешированный пароль (здесь он и хешируется)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error)
}

type userService struct {
	userRep userrep.UserRep
}

func NewUserService(empRep userrep.UserRep) UserService {
	return &userService{
		userRep: empRep,
	}
}

func (e *userService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return e.userRep.GetAll(ctx)
}

func (e *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return e.userRep.GetByID(ctx, id)
}

// func (e *userService) GetByLogin(login string) (*models.User, error) {
// 	return e.userRep.GetByLogin(login)
// }

// func (e *userService) Add(emp *models.User) error {
// 	return e.userRep.Add(emp)
// }

func (e *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return e.userRep.Delete(ctx, id)
}

func (e *userService) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error) {
	return e.userRep.Update(ctx, id, funcUpdate)
}
