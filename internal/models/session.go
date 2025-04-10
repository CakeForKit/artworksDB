package models

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	id        uuid.UUID
	iserID    uuid.UUID
	createdAt time.Time
}

// func NewUserSession()
