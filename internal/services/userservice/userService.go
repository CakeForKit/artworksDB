package userservice

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
)

type UserService interface {
	ChangeSubscribeToMailing(ctx context.Context, subscr bool) error
	GetSelf(ctx context.Context) (*models.User, error)
}

type userService struct {
	userRep userrep.UserRep
	authZ   auth.AuthZ
}

func NewUserService(userRep userrep.UserRep, authZ auth.AuthZ) UserService {
	return &userService{
		userRep: userRep,
		authZ:   authZ,
	}
}

func (m *userService) ChangeSubscribeToMailing(ctx context.Context, subscr bool) error {
	userID, err := m.authZ.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("userService.ChangeSubscribeToMailing: %w", err)
	}
	return m.userRep.UpdateSubscribeToMailing(ctx, userID, subscr)
}

func (m *userService) GetSelf(ctx context.Context) (*models.User, error) {
	userID, err := m.authZ.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("userService.GetSelf: %w", err)
	}
	user, err := m.userRep.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("userService.GetSelf: %w", err)
	}
	return user, nil
}
