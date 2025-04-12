package models

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id             uuid.UUID
	username       string
	login          string // unique
	hashedPassword string
	createdAt      time.Time
	mail           string
	subscribeMail  bool
}

var (
	ErrUserEmptyUsername    = errors.New("empty username")
	ErrUserEmptyLogin       = errors.New("empty login")
	ErrUserEmptyPassword    = errors.New("empty password")
	ErrUserInvalidCreatedAt = errors.New("invalid createdAt time")
	ErrUserInvalidEmail     = errors.New("invalid email format")
)

func NewUser(id uuid.UUID, username string, login string, hashedPassword string, createdAt time.Time, mail string, subscribeMail bool) (User, error) {
	if username == "" {
		return User{}, ErrUserEmptyUsername
	} else if login == "" {
		return User{}, ErrUserEmptyLogin
	} else if hashedPassword == "" {
		return User{}, ErrUserEmptyPassword
	} else if createdAt.IsZero() {
		return User{}, ErrUserInvalidCreatedAt
	} else if !isValidEmail(mail) {
		return User{}, ErrUserInvalidEmail
	}

	return User{
		id:             id,
		username:       username,
		login:          login,
		hashedPassword: hashedPassword,
		createdAt:      createdAt,
		mail:           mail,
		subscribeMail:  subscribeMail,
	}, nil
}

// GetID возвращает идентификатор пользователя
func (u *User) GetID() uuid.UUID {
	return u.id
}

// GetUsername возвращает имя пользователя
func (u *User) GetUsername() string {
	return u.username
}

// GetLogin возвращает логин пользователя
func (u *User) GetLogin() string {
	return u.login
}

// GetHashedPassword возвращает хэшированный пароль пользователя
func (u *User) GetHashedPassword() string {
	return u.hashedPassword
}

// GetCreatedAt возвращает дату создания пользователя
func (u *User) GetCreatedAt() time.Time {
	return u.createdAt
}

// GetMail возвращает email пользователя
func (u *User) GetMail() string {
	return u.mail
}

// IsSubscribedToMail возвращает статус подписки на рассылку
func (u *User) IsSubscribedToMail() bool {
	return u.subscribeMail
}

// Вспомогательная функция для проверки email
// func isValidEmail(email string) bool {
// 	fmt.Print(len(email) >= 3, contains(email, "@"), contains(email, "."))
// 	return len(email) >= 3 && contains(email, "@") && contains(email, ".")
// }

// func contains(s, substr string) bool {
// 	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
// }

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return (err == nil)
}
