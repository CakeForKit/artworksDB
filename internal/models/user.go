package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	id             uuid.UUID
	username       string
	login          string
	hashedPassword string
	createdAt      time.Time
	mail           string
	subscribe_mail bool
}
