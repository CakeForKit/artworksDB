package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	id             uuid.UUID
	username       string
	login          string // unique
	hashedPassword string
	createdAt      time.Time
	valid          bool
}

var (
	ErrAdminEmptyUsername    = errors.New("empty username")
	ErrAdminEmptyLogin       = errors.New("empty login")
	ErrAdminEmptyPassword    = errors.New("empty password")
	ErrAdminInvalidCreatedAt = errors.New("invalid createdAt time")
)

func NewAdmin(id uuid.UUID, username string, login string, hashedPassword string, createdAt time.Time, valid bool) (Admin, error) {
	admin := Admin{
		id:             id,
		username:       strings.TrimSpace(username),
		login:          strings.TrimSpace(login),
		hashedPassword: hashedPassword,
		createdAt:      createdAt,
		valid:          valid,
	}
	err := admin.validate()
	if err != nil {
		return Admin{}, err
	}
	return admin, nil
}

func (a *Admin) validate() error {
	if a.username == "" || len(a.username) > 50 {
		return ErrAdminEmptyUsername
	} else if a.login == "" || len(a.login) > 50 {
		return ErrAdminEmptyLogin
	} else if a.hashedPassword == "" || len(a.hashedPassword) > 255 {
		return ErrAdminEmptyPassword
	} else if a.createdAt.IsZero() {
		return ErrAdminInvalidCreatedAt
	}
	return nil
}

// GetID возвращает идентификатор администратора
func (a *Admin) GetID() uuid.UUID {
	return a.id
}

// GetUsername возвращает имя администратора
func (a *Admin) GetUsername() string {
	return a.username
}

// GetLogin возвращает логин администратора
func (a *Admin) GetLogin() string {
	return a.login
}

// GetHashedPassword возвращает хэшированный пароль администратора
func (a *Admin) GetHashedPassword() string {
	return a.hashedPassword
}

// GetCreatedAt возвращает дату создания учетной записи администратора
func (a *Admin) GetCreatedAt() time.Time {
	return a.createdAt
}

// IsValid возвращает статус активности учетной записи администратора
func (a *Admin) IsValid() bool {
	return a.valid
}

// SetValid устанавливает статус активности учетной записи администратора
func (a *Admin) SetValid(valid bool) {
	a.valid = valid
}
