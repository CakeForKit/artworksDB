package auth

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/google/uuid"
)

type authZContextKey int

const (
	AuthZContextKey authZContextKey = iota
)

var (
	ErrNotAuthZ    = errors.New("not authorized")
	ErrHasNoRights = errors.New("has no rights")
)

type AuthZ interface {
	Authorize(ctx context.Context, payload token.Payload) context.Context
	UserIDFromContext(ctx context.Context) (uuid.UUID, error)
	EmployeeIDFromContext(ctx context.Context) (uuid.UUID, error)
	AdminIDFromContext(ctx context.Context) (uuid.UUID, error)
}

func NewAuthZ() (AuthZ, error) {
	return &authZ{}, nil
}

type authZ struct {
}

func (a *authZ) Authorize(ctx context.Context, payload token.Payload) context.Context {
	return context.WithValue(ctx, AuthZContextKey, payload)
}

func (a *authZ) UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	payload, ok := ctx.Value(AuthZContextKey).(token.Payload)
	if !ok {
		return uuid.Nil, ErrNotAuthZ
	}
	if payload.Role != token.UserRole {
		return uuid.Nil, ErrHasNoRights
	}
	return payload.PersonID, nil
}

func (a *authZ) EmployeeIDFromContext(ctx context.Context) (uuid.UUID, error) {
	payload, ok := ctx.Value(AuthZContextKey).(token.Payload)
	if !ok {
		return uuid.Nil, ErrNotAuthZ
	}
	if payload.Role != token.EmployeeRole {
		return uuid.Nil, ErrHasNoRights
	}
	return payload.PersonID, nil
}

func (a *authZ) AdminIDFromContext(ctx context.Context) (uuid.UUID, error) {
	payload, ok := ctx.Value(AuthZContextKey).(token.Payload)
	if !ok {
		return uuid.Nil, ErrNotAuthZ
	}
	if payload.Role != token.EmployeeRole && payload.Role != token.AdminRole {
		return uuid.Nil, ErrHasNoRights
	}
	return payload.PersonID, nil
}
