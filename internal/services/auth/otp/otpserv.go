package otp

import (
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
)

type OTPService interface {
	SendOTP(user models.User) (string, error) // отправить код на почту (или др)
	// CheckOTP()
}

type EmailOTPServ struct {
}

func NewOTPService() OTPService {
	return &EmailOTPServ{}
}

func (s *EmailOTPServ) SendOTP(user models.User) (string, error) {
	return uuid.New().String()[:4], nil
}
