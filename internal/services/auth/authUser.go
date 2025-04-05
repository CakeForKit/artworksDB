package auth

import (
	"fmt"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/util"
	"github.com/google/uuid"
)

type LoginUserRequest struct {
	Login    string
	Password string
}

type RegisterUserRequest struct {
	Username      string
	Login         string
	Password      string
	Mail          string
	SubscribeMail bool
}

type AuthUserService interface {
	LoginUser(lur LoginUserRequest) (string, error)
	RegisterUser(rur RegisterUserRequest) error
}

func NewAuthUserService(config util.Config, urep userrep.UserRep) (AuthUserService, error) {
	tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	hasher, err := hasher.NewHasher()
	if err != nil {
		return nil, err
	}

	server := &authUserService{
		tokenMaker: tokenMaker,
		config:     config,
		userrep:    urep,
		hasher:     hasher,
	}

	return server, nil
}

type authUserService struct {
	tokenMaker token.TokenMaker
	config     util.Config
	userrep    userrep.UserRep
	hasher     hasher.Hasher
}

func (s *authUserService) LoginUser(lur LoginUserRequest) (string, error) {
	user, err := s.userrep.GetByLogin((lur.Login))
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(lur.Password, user.GetHashedPassword())
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.GetID(),
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authUserService) RegisterUser(rur RegisterUserRequest) error {
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
		rur.Mail,
		rur.SubscribeMail,
	)
	if err != nil {
		return nil
	}
	err = s.userrep.Add(&user)
	return err
}
