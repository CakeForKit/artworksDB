package testobj

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type UserMother interface {
	DefaultUser(userID uuid.UUID) models.User
	UserWithPswdHash(userID uuid.UUID, hashedPassword string) models.User
	DefaultUserP(userID uuid.UUID) *models.User
}

func NewUserMother() UserMother {
	return &userMother{}
}

type userMother struct{}

func (um *userMother) DefaultUser(userID uuid.UUID) models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		"test-login",
		"hashed-password",
		time.Now(),
		"user@test.com",
		true,
	)
	return user
}

func (um *userMother) UserWithPswdHash(userID uuid.UUID, hashedPassword string) models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		"test-login",
		hashedPassword,
		time.Now(),
		"user@test.com",
		true,
	)
	return user
}

func (um *userMother) DefaultUserP(userID uuid.UUID) *models.User {
	user, _ := models.NewUser(
		userID,
		"test-user",
		"test-login",
		"hashed-password",
		time.Now(),
		"user@test.com",
		true,
	)
	return &user
}
