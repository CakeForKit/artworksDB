package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	// duration - корректный срок действия, return - одписанную строку токена или ошибку
	CreateToken(id uuid.UUID, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)
