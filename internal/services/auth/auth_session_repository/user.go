package authsessionrep

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Сессия otp (между тем как user отправил логин и пароль и подтвердил код)
type AuthUserSessionRep interface {
	Create(UserID uuid.UUID, RequiresOTP string) uuid.UUID
	CheckOTP(sessionID uuid.UUID, otp string) (*OTPSession, error)
	Remove(sessionID uuid.UUID)
}

var (
	ErrAuthUserSessionRep = errors.New("AuthUserSessionRep")
	ErrSessionNotFound    = errors.New("OTP session not found")
	ErrSerialyzeSession   = errors.New("error serialyze session value")
	ErrExpiresSession     = errors.New("expires session")
	ErrWrongOTP           = errors.New("wrong otp")
)

type OTPSession struct {
	SessionID   uuid.UUID
	UserID      uuid.UUID
	RequiresOTP string
	CreatedAt   time.Time
}

type authUserSessionRep struct {
	sessions        sync.Map
	durationSession time.Duration
}

func NewAuthUserSessionRep(durationSession time.Duration) AuthUserSessionRep {
	return &authUserSessionRep{
		sessions:        sync.Map{},
		durationSession: durationSession,
	}
}

func (r *authUserSessionRep) Create(UserID uuid.UUID, RequiresOTP string) uuid.UUID {
	s := uuid.New()
	r.sessions.Store(s, &OTPSession{
		SessionID:   s,
		UserID:      UserID,
		RequiresOTP: RequiresOTP,
		CreatedAt:   time.Now().UTC(),
	})
	return s
}

func (r *authUserSessionRep) CheckOTP(sessionID uuid.UUID, otp string) (*OTPSession, error) {
	baseErr := fmt.Errorf("%w: CheckOTP", ErrAuthUserSessionRep)
	val, ok := r.sessions.Load(sessionID)
	if !ok {
		return nil, fmt.Errorf("%w: %w", baseErr, ErrSessionNotFound)
	}
	session, ok := val.(*OTPSession)
	if !ok {
		return nil, fmt.Errorf("%w: %w", baseErr, ErrSerialyzeSession)
	}
	if time.Now().UTC().Sub(session.CreatedAt) > r.durationSession {
		return nil, fmt.Errorf("%w: %w", baseErr, ErrExpiresSession)
	}
	if session.RequiresOTP != otp {
		return nil, fmt.Errorf("%w: %w", baseErr, ErrWrongOTP)
	}
	return session, nil
}

func (r *authUserSessionRep) Remove(sessionID uuid.UUID) {
	r.sessions.Delete(sessionID)
}
