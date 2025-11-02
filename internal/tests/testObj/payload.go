package testobj

import (
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/google/uuid"
)

type PayloadMother interface {
	UserPayload(userID uuid.UUID) token.Payload
	EmployeePayload(emplID uuid.UUID) token.Payload
	AdminPayload(adminID uuid.UUID) token.Payload
}

func NewPayloadMother() PayloadMother {
	return &payloadMother{}
}

type payloadMother struct{}

func (pm *payloadMother) UserPayload(userID uuid.UUID) token.Payload {
	return token.Payload{
		PersonID:  userID,
		Role:      token.UserRole,
		ExpiredAt: time.Now().Add(time.Hour),
	}
}

func (pm *payloadMother) EmployeePayload(emplID uuid.UUID) token.Payload {
	return token.Payload{
		PersonID:  emplID,
		Role:      token.EmployeeRole,
		ExpiredAt: time.Now().Add(time.Hour),
	}
}

func (pm *payloadMother) AdminPayload(adminID uuid.UUID) token.Payload {
	return token.Payload{
		PersonID:  adminID,
		Role:      token.AdminRole,
		ExpiredAt: time.Now().Add(time.Hour),
	}
}
