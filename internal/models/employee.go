package models

import (
	"errors"
	"strings"
	"time"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type Employee struct {
	id             uuid.UUID
	username       string
	login          string
	hashedPassword string
	createdAt      time.Time
	valid          bool
	adminID        uuid.UUID
}

var (
	ErrEmployeeEmptyUsername    = errors.New("empty username")
	ErrEmployeeEmptyLogin       = errors.New("empty login")
	ErrEmployeeEmptyPassword    = errors.New("empty password")
	ErrEmployeeInvalidCreatedAt = errors.New("invalid createdAt time")
	ErrEmployeeInvalidAdminID   = errors.New("invalid admin ID")
)

func NewEmployee(id uuid.UUID, username string, login string, hashedPassword string,
	createdAt time.Time, valid bool, adminID uuid.UUID) (Employee, error) {
	employee := Employee{
		id:             id,
		username:       strings.TrimSpace(username),
		login:          strings.TrimSpace(login),
		hashedPassword: hashedPassword,
		createdAt:      createdAt,
		valid:          valid,
		adminID:        adminID,
	}
	err := employee.validate()
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}

func (e *Employee) validate() error {
	if e.username == "" || len(e.username) > 50 {
		return ErrEmployeeEmptyUsername
	} else if e.login == "" || len(e.login) > 50 {
		return ErrEmployeeEmptyLogin
	} else if e.hashedPassword == "" || len(e.hashedPassword) > 255 {
		return ErrEmployeeEmptyPassword
	} else if e.createdAt.IsZero() {
		return ErrEmployeeInvalidCreatedAt
	} else if e.adminID == uuid.Nil {
		return ErrEmployeeInvalidAdminID
	}
	return nil
}

func (e *Employee) GetID() uuid.UUID {
	return e.id
}

func (e *Employee) GetUsername() string {
	return e.username
}

func (e *Employee) GetLogin() string {
	return e.login
}

func (e *Employee) GetHashedPassword() string {
	return e.hashedPassword
}

func (e *Employee) GetCreatedAt() time.Time {
	return e.createdAt
}

func (e *Employee) IsValid() bool {
	return e.valid
}

func (e *Employee) SetValid(valid bool) {
	e.valid = valid
}

func (e *Employee) GetAdminID() uuid.UUID {
	return e.adminID
}

func (e *Employee) ToEmployeeResponse() jsonreqresp.EmployeeResponse {
	return jsonreqresp.EmployeeResponse{
		ID:        e.id.String(),
		Username:  e.username,
		Login:     e.login,
		CreatedAt: e.createdAt,
		Valid:     e.valid,
		AdminID:   e.adminID.String(),
	}
}
