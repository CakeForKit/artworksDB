package models

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id             uuid.UUID
	username       string
	login          string // unique
	hashedPassword string
	createdAt      time.Time
	email          string
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
	user := User{
		id:             id,
		username:       strings.TrimSpace(username),
		login:          strings.TrimSpace(login),
		hashedPassword: hashedPassword,
		createdAt:      createdAt,
		email:          strings.TrimSpace(mail),
		subscribeMail:  subscribeMail,
	}
	err := user.validate()
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *User) validate() error {
	if u.username == "" || len(u.username) > 50 {
		return ErrUserEmptyUsername
	} else if u.login == "" || len(u.login) > 50 {
		return ErrUserEmptyLogin
	} else if u.hashedPassword == "" || len(u.hashedPassword) > 255 {
		return ErrUserEmptyPassword
	} else if u.createdAt.IsZero() {
		return ErrUserInvalidCreatedAt
	} else if len(u.email) > 100 || !isValidEmail(u.email) {
		return ErrUserInvalidEmail
	}
	return nil
}

func CmpUsers(u1, u2 *User) bool {
	if u1 == nil || u2 == nil {
		return u1 == u2
	}
	// createdAt - не сравниваем
	return u1.id == u2.id &&
		u1.username == u2.username &&
		u1.login == u2.login &&
		u1.hashedPassword == u2.hashedPassword &&
		u1.email == u2.email &&
		u1.subscribeMail == u2.subscribeMail
}

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetLogin() string {
	return u.login
}

func (u *User) GetHashedPassword() string {
	return u.hashedPassword
}

func (u *User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *User) GetEmail() string {
	return u.email
}

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
