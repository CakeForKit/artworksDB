package otp

import (
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockOTPService struct {
	mock.Mock
}

func (m *MockOTPService) SendOTP(user models.User) (string, error) {
	args := m.Called(user)
	return args.Get(0).(string), args.Error(1)
}
