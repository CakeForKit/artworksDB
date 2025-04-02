package hasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(password string) (string, error) // HashPassword возвращает bcrypt хэш пароля
	CheckPassword(password string, hashedPassword string) error
}

var (
	ErrHash     = errors.New("failed to hash password")
	ErrPassword = errors.New("error CompareHashAndPassword")
)

func NewHasher() (Hasher, error) {
	return &bcryptHasher{}, nil
}

// ----------------

type bcryptHasher struct {
}

func (h *bcryptHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHash
	}
	return string(hashedPassword), nil
}

func (h *bcryptHasher) CheckPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		err = ErrPassword
	}
	return err
}
