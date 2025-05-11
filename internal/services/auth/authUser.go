package auth

import (
	"context"
	"fmt"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/google/uuid"
)

type LoginUserRequest struct {
	Login    string `json:"login" binding:"required,alphanum,min=4,max=50" example:"ulogin"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum,max=50" example:"uname"`
	Login          string `json:"login" binding:"required,alphanum,min=4,max=50" example:"ulogin"`
	Password       string `json:"password" binding:"required,min=4" example:"12345678"`
	Email          string `json:"email" binding:"required,email,min=6,max=100" example:"uuser@test.ru"`
	SubscribeEmail bool   `json:"subscribe_email" binding:"required,boolean" example:"true"`
}

type AuthUser interface {
	LoginUser(ctx context.Context, lur LoginUserRequest) (string, error)
	RegisterUser(ctx context.Context, rur RegisterUserRequest) error
	VerifyByToken(token string) (*token.Payload, error)
}

type authUser struct {
	tokenMaker token.TokenMaker
	config     cnfg.AppConfig
	userrep    userrep.UserRep
	hasher     hasher.Hasher
}

func NewAuthUser(config cnfg.AppConfig, urep userrep.UserRep) (AuthUser, error) {
	tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	hasher, err := hasher.NewHasher()
	if err != nil {
		return nil, err
	}

	server := &authUser{
		tokenMaker: tokenMaker,
		config:     config,
		userrep:    urep,
		hasher:     hasher,
	}

	return server, nil
}

func (s *authUser) LoginUser(ctx context.Context, lur LoginUserRequest) (string, error) {
	user, err := s.userrep.GetByLogin(ctx, lur.Login)
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(lur.Password, user.GetHashedPassword())
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.GetID(),
		token.UserRole,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authUser) RegisterUser(ctx context.Context, rur RegisterUserRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return err
	}
	user, err := models.NewUser(
		uuid.New(),
		rur.Username,
		rur.Login,
		hashedPassword,
		time.Now(),
		rur.Email,
		rur.SubscribeEmail,
	)
	if err != nil {
		return nil
	}
	err = s.userrep.Add(ctx, &user)
	return err
}

func (s *authUser) VerifyByToken(tokenStr string) (*token.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, token.UserRole)
}
