package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	attemptsrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/attempts_rep"
	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
	authsessionrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_session_repository"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/otp"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/google/uuid"
)

type AuthUser interface {
	RegisterUser(ctx context.Context, rur authmodels.RegisterUserRequest) error
	LoginUser(ctx context.Context, lur authmodels.LoginUserRequest) (uuid.UUID, error)
	OTP(ctx context.Context, sessionID uuid.UUID, otp string) (string, error)
	VerifyByToken(token string) (*token.Payload, error)
}

var (
	ErrAuthUser = errors.New("AuthUser")
)

type authUser struct {
	tokenMaker      token.TokenMaker
	config          cnfg.AppConfig
	userrep         userrep.UserRep
	hasher          hasher.Hasher
	otpServ         otp.OTPService
	sessionRep      authsessionrep.AuthUserSessionRep
	loginAttemptRep attemptsrep.LoginAttemptUserRep
	otpAttemptRep   attemptsrep.OTPAttemptRep
}

func NewAuthUser(
	config cnfg.AppConfig, urep userrep.UserRep,
	tokenMaker token.TokenMaker, hasher hasher.Hasher,
	otpServ otp.OTPService,
	authUserSessionRep authsessionrep.AuthUserSessionRep,
	loginAttemptRep attemptsrep.LoginAttemptUserRep,
	otpAttemptRep attemptsrep.OTPAttemptRep,
) (AuthUser, error) {
	server := &authUser{
		tokenMaker:      tokenMaker,
		config:          config,
		userrep:         urep,
		hasher:          hasher,
		otpServ:         otpServ,
		sessionRep:      authUserSessionRep,
		loginAttemptRep: loginAttemptRep,
		otpAttemptRep:   otpAttemptRep,
	}

	return server, nil
}

func (s *authUser) LoginUser(ctx context.Context, lur authmodels.LoginUserRequest) (uuid.UUID, error) {
	if err := s.loginAttemptRep.Add(lur.Login); err != nil {
		return uuid.Nil, fmt.Errorf("%v: %w", ErrAuthUser, err)
	}
	user, err := s.userrep.GetByLogin(ctx, lur.Login)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%v: %w", ErrAuthUser, err)
	}

	err = s.hasher.CheckPassword(lur.Password, user.GetHashedPassword())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%v: %w", ErrAuthUser, err)
	}

	requiresOTP, err := s.otpServ.SendOTP(*user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%v: %w", ErrAuthUser, err)
	}
	sessionID := s.sessionRep.Create(user.GetID(), requiresOTP)

	s.loginAttemptRep.Remove(lur.Login)
	return sessionID, nil
}

func (s *authUser) OTP(ctx context.Context, sessionID uuid.UUID, otp string) (string, error) {
	if err := s.otpAttemptRep.Add(sessionID); err != nil {
		return "", fmt.Errorf("%v: %w", ErrAuthUser, err)
	}

	session, err := s.sessionRep.CheckOTP(sessionID, otp)
	if err != nil {
		return "", fmt.Errorf("%v: %w", ErrAuthUser, err)
	}

	accessToken, err := s.tokenMaker.CreateToken(
		session.UserID,
		token.UserRole,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}

	s.otpAttemptRep.Remove(sessionID)
	return accessToken, nil
}

func (s *authUser) RegisterUser(ctx context.Context, rur authmodels.RegisterUserRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return err
	}
	user, err := models.NewUser(
		uuid.New(),
		rur.Username,
		rur.Login,
		hashedPassword,
		time.Now(),
		rur.Email,
		rur.SubscribeEmail,
	)
	if err != nil {
		return nil
	}
	err = s.userrep.Add(ctx, &user)
	return err
}

func (s *authUser) VerifyByToken(tokenStr string) (*token.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, token.UserRole)
}
