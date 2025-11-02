package otp

import (
	"errors"
	"fmt"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailserv"
	"github.com/google/uuid"
)

type OTPService interface {
	SendOTP(user models.User) (string, error) // отправить код на почту (или др)
}

var (
	ErrOTPService = errors.New("OTPService")
)

type EmailOTPServ struct {
	emailServ emailserv.EmailService
}

func NewOTPService() OTPService {
	return &EmailOTPServ{}
}

func (s *EmailOTPServ) SendOTP(user models.User) (string, error) {
	otpCode := uuid.New().String()[:4]
	err := s.emailServ.SendEmail([]string{user.GetEmail()}, "Код подтверждения", otpCode)
	if err != nil {
		return "", fmt.Errorf("%v: %w", ErrOTPService, err)
	}
	return otpCode, nil
}
