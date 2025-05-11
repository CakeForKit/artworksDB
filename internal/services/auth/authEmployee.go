package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/google/uuid"
)

type LoginEmployeeRequest struct {
	Login    string `json:"login" binding:"required,alphanum,min=4,max=50" example:"elogin"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type RegisterEmployeeRequest struct {
	Username string `json:"username" binding:"required,alphanum,max=50" example:"ename"`
	Login    string `json:"login" binding:"required,alphanum,min=4,max=50" example:"elogin"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
	// Valid    bool      `json:"valid" binding:"required,boolean" example:"true"`
	// AdminID uuid.UUID `json:"adminID" binding:"required,uuid" example:"8f005053-5b95-4a6a-bdcd-7395ee3ed204"`
}

type LoginEmployeeResponse struct {
	AccessToken string `json:"access_token"`
}

type AuthEmployee interface {
	LoginEmployee(ctx context.Context, ler LoginEmployeeRequest) (string, error)
	RegisterEmployee(ctx context.Context, rer RegisterEmployeeRequest, adminID uuid.UUID) error
	VerifyByToken(tokenStr string) (*token.Payload, error)
}

var (
	ErrEmployeeNotValid = errors.New("the Employee has no rights")
)

type authEmployee struct {
	tokenMaker  token.TokenMaker
	config      cnfg.AppConfig
	employeerep employeerep.EmployeeRep
	hasher      hasher.Hasher
}

func NewAuthEmployee(config cnfg.AppConfig, erep employeerep.EmployeeRep) (AuthEmployee, error) {
	tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	hasher, err := hasher.NewHasher()
	if err != nil {
		return nil, err
	}

	service := &authEmployee{
		tokenMaker:  tokenMaker,
		config:      config,
		hasher:      hasher,
		employeerep: erep,
	}

	return service, nil
}

func (s *authEmployee) LoginEmployee(ctx context.Context, ler LoginEmployeeRequest) (string, error) {
	employee, err := s.employeerep.GetByLogin(ctx, ler.Login)
	if err != nil {
		return "", err
	}

	if !employee.IsValid() {
		return "", ErrEmployeeNotValid
	}

	err = s.hasher.CheckPassword(ler.Password, employee.GetHashedPassword())
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		employee.GetID(),
		token.EmployeeRole,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authEmployee) RegisterEmployee(ctx context.Context, rer RegisterEmployeeRequest, adminID uuid.UUID) error {
	hashedPassword, err := s.hasher.HashPassword(rer.Password)
	if err != nil {
		return err
	}
	employee, err := models.NewEmployee(
		uuid.New(),
		rer.Username,
		rer.Login,
		hashedPassword,
		time.Now(),
		true,
		adminID,
	)
	if err != nil {
		return err
	}
	err = s.employeerep.Add(ctx, &employee)
	return err
}

func (s *authEmployee) VerifyByToken(tokenStr string) (*token.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, token.EmployeeRole)
}
