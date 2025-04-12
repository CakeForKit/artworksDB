package userservice

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers() []*models.User
	GetByID(id uuid.UUID) (*models.User, error)
	// GetByLogin(login string) (*models.User, error)
	// Add(*models.User) error // нехешированный пароль (здесь он и хешируется)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error)
}

type userService struct {
	userRep userrep.UserRep
}

func NewUserService(empRep userrep.UserRep) UserService {
	return &userService{
		userRep: empRep,
	}
}

func (e *userService) GetAllUsers() []*models.User {
	return e.userRep.GetAll()
}

func (e *userService) GetByID(id uuid.UUID) (*models.User, error) {
	return e.userRep.GetByID(id)
}

// func (e *userService) GetByLogin(login string) (*models.User, error) {
// 	return e.userRep.GetByLogin(login)
// }

// func (e *userService) Add(emp *models.User) error {
// 	return e.userRep.Add(emp)
// }

func (e *userService) Delete(id uuid.UUID) error {
	return e.userRep.Delete(id)
}

func (e *userService) Update(id uuid.UUID, funcUpdate func(*models.User) (*models.User, error)) (*models.User, error) {
	return e.userRep.Update(id, funcUpdate)
}
