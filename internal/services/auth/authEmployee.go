package auth

import (
	"context"
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
	Login    string
	Password string
}

type RegisterEmployeeRequest struct {
	Username string
	Login    string
	Password string
	Valid    bool
	AdminID  uuid.UUID
}

type AuthEmployee interface {
	LoginEmployee(ctx context.Context, ler LoginEmployeeRequest) (string, error)
	RegisterEmployee(ctx context.Context, rer RegisterEmployeeRequest) error
}

func NewAuthEmployee(config cnfg.AppConfig) (AuthEmployee, error) {
	tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	service := &authEmployee{
		tokenMaker: tokenMaker,
		config:     config,
	}

	return service, nil
}

type authEmployee struct {
	tokenMaker  token.TokenMaker
	config      cnfg.AppConfig
	employeerep employeerep.EmployeeRep
	hasher      hasher.Hasher
}

func (s *authEmployee) LoginEmployee(ctx context.Context, ler LoginEmployeeRequest) (string, error) {
	employee, err := s.employeerep.GetByLogin(ctx, ler.Login)
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(ler.Password, employee.GetHashedPassword())
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		employee.GetID(),
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authEmployee) RegisterEmployee(ctx context.Context, rer RegisterEmployeeRequest) error {
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
		rer.Valid,
		rer.AdminID,
	)
	if err != nil {
		return err
	}
	err = s.employeerep.Add(ctx, &employee)
	return err
}
