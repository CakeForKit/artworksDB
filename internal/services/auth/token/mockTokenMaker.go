package token

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTokenMaker struct {
	mock.Mock
}

func (m *MockTokenMaker) CreateToken(userID uuid.UUID, role string, duration time.Duration) (string, error) {
	args := m.Called(userID, role, duration)
	return args.String(0), args.Error(1)
}

func (m *MockTokenMaker) VerifyToken(tokenStr string, expectedRole string) (*Payload, error) {
	args := m.Called(tokenStr, expectedRole)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Payload), args.Error(1)
}
