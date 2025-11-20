package authadmin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/google/uuid"
)

type LoginAdminRequest struct {
	Login    string `json:"login" binding:"required,alphanum,max=50" example:"admin"`
	Password string `json:"password" binding:"required,min=6" example:"12345678"`
}

type LoginAdminResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterAdminRequest struct {
	Adminname string `json:"adminname" binding:"required,alphanum,max=50" example:"admin"`
	Login     string `json:"login" binding:"required,alphanum,max=50" example:"admin"`
	Password  string `json:"password" binding:"required,min=6" example:"12345678"`
}

type AdminAuth interface {
	LoginAdmin(ctx context.Context, lur LoginAdminRequest) (string, error)
	RegisterAdmin(ctx context.Context, rur RegisterAdminRequest) error
	VerifyByToken(token string) (*token.Payload, error)
}

var (
	ErrAdminNotValid = errors.New("the Admin has no rights")
)

type authAdmin struct {
	tokenMaker token.TokenMaker
	config     cnfg.AppConfig
	adminrep   adminrep.AdminRep
	hasher     hasher.Hasher
}

func NewAuthAdmin(config cnfg.AppConfig, urep adminrep.AdminRep,
	tokenMaker token.TokenMaker, hasher hasher.Hasher) AdminAuth {
	server := &authAdmin{
		tokenMaker: tokenMaker,
		config:     config,
		adminrep:   urep,
		hasher:     hasher,
	}

	return server
}

func (s *authAdmin) LoginAdmin(ctx context.Context, lur LoginAdminRequest) (string, error) {
	admin, err := s.adminrep.GetByLogin(ctx, lur.Login)
	if err != nil {
		return "", fmt.Errorf("LoginAdmin: %w", err)
	}

	if !admin.IsValid() {
		return "", ErrAdminNotValid
	}

	err = s.hasher.CheckPassword(lur.Password, admin.GetHashedPassword())
	if err != nil {
		return "", fmt.Errorf("LoginAdmin: %w", err)
	}

	accessToken, err := s.tokenMaker.CreateToken(
		admin.GetID(),
		token.AdminRole,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", fmt.Errorf("LoginAdmin: %w", err)
	}
	return accessToken, nil
}

func (s *authAdmin) RegisterAdmin(ctx context.Context, rur RegisterAdminRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return fmt.Errorf("LoginAdmin: %w", err)
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
	if err != nil {
		return fmt.Errorf("LoginAdmin: %w", err)
	}
	return nil
}

func (s *authAdmin) VerifyByToken(tokenStr string) (*token.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, token.AdminRole)
}
