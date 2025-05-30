package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TokenMaker interface {
	// duration - корректный срок действия, return - одписанную строку токена или ошибку
	CreateToken(id uuid.UUID, role string, duration time.Duration) (string, error)
	VerifyToken(token string, role string) (*Payload, error)
}

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrExpiredToken  = errors.New("token has expired")
	ErrIncorrectRole = errors.New("incorrect role")
)

func NewTokenMaker(symmetricKey string) (TokenMaker, error) {
	return NewPasetoMaker(symmetricKey)
}
