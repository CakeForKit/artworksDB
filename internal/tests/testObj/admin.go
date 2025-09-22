package testobj

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type AdminMother interface {
	DefaultAdmin(adminID uuid.UUID) models.Admin
	AdminWithPswdHash(adminID uuid.UUID, hashedPassword string) models.Admin
}

func NewAdminMother() AdminMother {
	return &adminMother{}
}

type adminMother struct{}

func (um *adminMother) DefaultAdmin(adminID uuid.UUID) models.Admin {
	admin, _ := models.NewAdmin(
		adminID,
		"test-admin",
		"test-login"+uuid.NewString(),
		"hashed-password"+uuid.NewString(),
		time.Now(),
		true,
	)
	return admin
}

func (um *adminMother) AdminWithPswdHash(adminID uuid.UUID, hashedPassword string) models.Admin {
	admin, _ := models.NewAdmin(
		adminID,
		"test-admin",
		"test-login"+uuid.NewString(),
		"hashed-password"+uuid.NewString(),
		time.Now(),
		true,
	)
	return admin
}
