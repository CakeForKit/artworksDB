package userservice

import (
	"context"
	"fmt"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
)

type UserService interface {
	ChangeSubscribeToMailing(ctx context.Context, subscr bool) error
	GetSelf(ctx context.Context) (*models.User, error)
	DeleteSelf(ctx context.Context) error
	ChangePassword(ctx context.Context, newPassword string) error
}

type userService struct {
	userRep userrep.UserRep
	authZ   auth.AuthZ
	hash    hasher.Hasher
}

func NewUserService(userRep userrep.UserRep, authZ auth.AuthZ, hash hasher.Hasher) UserService {
	return &userService{
		userRep: userRep,
		authZ:   authZ,
		hash:    hash,
	}
}

func (m *userService) ChangeSubscribeToMailing(ctx context.Context, subscr bool) error {
	userID, err := m.authZ.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("userService.ChangeSubscribeToMailing: %w", err)
	}
	return m.userRep.UpdateSubscribeToMailing(ctx, userID, subscr)
}

func (m *userService) ChangePassword(ctx context.Context, newPassword string) error {
	userID, err := m.authZ.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("userService.ChangePassword: %w", err)
	}
	_, err = m.userRep.Update(ctx, userID, func(u *models.User) (*models.User, error) {
		h, err := m.hash.HashPassword(newPassword)
		if err != nil {
			return nil, err
		}
		updated, err := models.NewUser(
			u.GetID(), u.GetUsername(), u.GetLogin(),
			h, time.Now().UTC(), u.GetEmail(), u.IsSubscribedToMail(),
		)
		return &updated, err
	})
	return err
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

func (m *userService) DeleteSelf(ctx context.Context) error {
	userID, err := m.authZ.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("userService.GetSelf: %w", err)
	}
	user, err := m.userRep.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("userService.GetSelf: %w", err)
	}
	err = m.userRep.Delete(ctx, user.GetID())
	if err != nil {
		return fmt.Errorf("userService.GetSelf: %w", err)
	}
	return nil
}
