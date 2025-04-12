package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	id             uuid.UUID
	username       string
	login          string
	hashedPassword string
	createdAt      time.Time
}

var (
	ErrEmployeeEmptyUsername = errors.New("empty username")
	ErrEmployeeEmptyLogin    = errors.New("empty login")
	ErrEmployeeEmptyPassword = errors.New("empty password")
	ErrEmployeeCreatedAt     = errors.New("invalid createdAt time")
)

func NewEmployee(id uuid.UUID, username string, login string, hashedPassword string, createdAt time.Time) (Employee, error) {
	if username == "" {
		return Employee{}, ErrEmployeeEmptyUsername
	} else if login == "" {
		return Employee{}, ErrEmployeeEmptyLogin
	} else if hashedPassword == "" {
		return Employee{}, ErrEmployeeEmptyPassword
	} else if createdAt.IsZero() {
		return Employee{}, ErrEmployeeCreatedAt
	}
	return Employee{
		id:             id,
		username:       username,
		login:          login,
		hashedPassword: hashedPassword,
		createdAt:      createdAt,
	}, nil
}

// GetID возвращает идентификатор сотрудника
func (e *Employee) GetID() uuid.UUID {
	return e.id
}

// GetUsername возвращает имя пользователя сотрудника
func (e *Employee) GetUsername() string {
	return e.username
}

// GetLogin возвращает логин сотрудника
func (e *Employee) GetLogin() string {
	return e.login
}

// GetHashedPassword возвращает хэшированный пароль сотрудника
func (e *Employee) GetHashedPassword() string {
	return e.hashedPassword
}
