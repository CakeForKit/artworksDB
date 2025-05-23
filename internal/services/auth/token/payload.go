package token

import (
	"time"

	"github.com/google/uuid"
)

// данные полезной нагрузки, хранящиеся внутри тела токена

const (
	UserRole     = "user_role"
	EmployeeRole = "employee_role"
	AdminRole    = "admin_role"
)

type Payload struct {
	// ID        uuid.UUID `json:"id"`
	PersonID uuid.UUID `json:"person_id"`
	Role     string    `json:"role"`
	// IssuedAt  time.Time `json:"issued_at"`  // время создания токена
	ExpiredAt time.Time `json:"expired_at"` // время когда срок действия токена истечет
}

func NewPayload(personID uuid.UUID, role string, duration time.Duration) (*Payload, error) {
	// tokenID, err := uuid.NewRandom()
	// if err != nil {
	// 	return nil, err
	// }

	payload := &Payload{
		// ID:        tokenID,
		PersonID: personID,
		Role:     role,
		// IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid(role string) error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	} else if payload.Role != role {
		return ErrIncorrectRole
	}
	return nil
}

// func (p *Payload) GetID() uuid.UUID {
// 	return p.ID
// }

func (p *Payload) GetPersonID() uuid.UUID {
	return p.PersonID
}

func (p *Payload) GetRole() string {
	return p.Role
}

// func (p *Payload) GetIssuedAt() time.Time {
// 	return p.IssuedAt
// }

func (p *Payload) GetExpiredAt() time.Time {
	return p.ExpiredAt
}
