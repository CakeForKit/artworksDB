package testobj

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type AdminMother interface {
	DefaultAdmin(adminID uuid.UUID) models.Admin
	DefaultAdminP(adminID uuid.UUID) *models.Admin
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
		"test-login"+adminID.String(),
		"hashed-password"+adminID.String(),
		time.Now(),
		true,
	)
	return admin
}

func (um *adminMother) DefaultAdminP(adminID uuid.UUID) *models.Admin {
	a := um.DefaultAdmin(adminID)
	return &a
}

func (um *adminMother) AdminWithPswdHash(adminID uuid.UUID, hashedPassword string) models.Admin {
	admin, _ := models.NewAdmin(
		adminID,
		"test-admin",
		"test-login"+adminID.String(),
		hashedPassword,
		time.Now().UTC(),
		true,
	)
	return admin
}
