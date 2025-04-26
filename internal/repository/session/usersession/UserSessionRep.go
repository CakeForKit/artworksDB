package usersession

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type UserSessionRep interface {
	GetByID(sid uuid.UUID) (*models.UserSession, error)
	Add(s *models.UserSession) error
	Delete(id uuid.UUID) error
	ClearExpired() error
}
