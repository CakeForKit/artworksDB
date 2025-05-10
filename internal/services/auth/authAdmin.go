package auth

import (
	"context"
	"fmt"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/google/uuid"
)

type LoginAdminRequest struct {
	Login    string `json:"login" binding:"required,alphanum,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginAdminResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterAdminRequest struct {
	Adminname string `json:"adminname" binding:"required,alphanum,max=50"`
	Login     string `json:"login" binding:"required,alphanum,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
}

type AuthAdmin interface {
	LoginAdmin(ctx context.Context, lur LoginAdminRequest) (string, error)
	RegisterAdmin(ctx context.Context, rur RegisterAdminRequest) error
	VerifyByToken(token string) (*token.Payload, error)
}

type authAdmin struct {
	tokenMaker token.TokenMaker
	config     cnfg.AppConfig
	adminrep   adminrep.AdminRep
	hasher     hasher.Hasher
}

func NewAuthAdmin(config cnfg.AppConfig, urep adminrep.AdminRep) (AuthAdmin, error) {
	tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	hasher, err := hasher.NewHasher()
	if err != nil {
		return nil, err
	}

	server := &authAdmin{
		tokenMaker: tokenMaker,
		config:     config,
		adminrep:   urep,
		hasher:     hasher,
	}

	return server, nil
}

func (s *authAdmin) LoginAdmin(ctx context.Context, lur LoginAdminRequest) (string, error) {
	admin, err := s.adminrep.GetByLogin(ctx, lur.Login)
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(lur.Password, admin.GetHashedPassword())
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		admin.GetID(),
		token.AdminRole,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authAdmin) RegisterAdmin(ctx context.Context, rur RegisterAdminRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return err
	}
	admin, err := models.NewAdmin(
		uuid.New(),
		rur.Adminname,
		rur.Login,
		hashedPassword,
		time.Now(),
		true,
	)
	if err != nil {
		return nil
	}
	err = s.adminrep.Add(ctx, &admin)
	return err
}

func (s *authAdmin) VerifyByToken(tokenStr string) (*token.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, token.AdminRole)
}
