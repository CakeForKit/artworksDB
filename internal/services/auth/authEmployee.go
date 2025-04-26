package auth

import (
	"fmt"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/config"
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
	LoginEmployee(ler LoginEmployeeRequest) (string, error)
	RegisterEmployee(rer RegisterEmployeeRequest) error
}

func NewAuthEmployee(config config.Config) (AuthEmployee, error) {
	tokenMaker, err := token.NewTokenMaker(config.App.TokenSymmetricKey)
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
	config      config.Config
	employeerep employeerep.EmployeeRep
	hasher      hasher.Hasher
}

func (s *authEmployee) LoginEmployee(ler LoginEmployeeRequest) (string, error) {
	employee, err := s.employeerep.GetByLogin(ler.Login)
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(ler.Password, employee.GetHashedPassword())
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		employee.GetID(),
		s.config.App.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authEmployee) RegisterEmployee(rer RegisterEmployeeRequest) error {
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
	err = s.employeerep.Add(&employee)
	return err
}
