package attemptsrep

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type otpAttempt struct {
	OTPSessionID uuid.UUID
	Count        int
	CreatedAt    time.Time
}

type OTPAttemptRep interface {
	Add(otpSessionID uuid.UUID) error
	Remove(otpSessionID uuid.UUID)
}

var (
	ErrOTPAttemptRep = errors.New("OTPAttemptRep")
)

type otpAttemptRep struct {
	maxOTPAttempts  int
	durationSession time.Duration
	attempts        sync.Map
}

func NewOTPAttemptRep(maxOTPAttempts int, durationSession time.Duration) OTPAttemptRep {
	return &otpAttemptRep{
		maxOTPAttempts:  maxOTPAttempts,
		durationSession: durationSession,
		attempts:        sync.Map{},
	}
}

func (r *otpAttemptRep) Add(otpSessionID uuid.UUID) error {
	baseErr := fmt.Errorf("%w: Add", ErrLoginAttemptUserRep)
	val, ok := r.attempts.Load(otpSessionID)
	if !ok {
		r.attempts.Store(otpSessionID, &otpAttempt{
			OTPSessionID: otpSessionID,
			Count:        1,
			CreatedAt:    time.Now().UTC(),
		})
		return nil
	}
	attempt, ok := val.(*otpAttempt)
	if !ok {
		return fmt.Errorf("%w: %w", baseErr, ErrSerialyzeSession)
	}
	if attempt.CreatedAt.Sub(time.Now().UTC()) > r.durationSession {
		attempt.Count = 1
		attempt.CreatedAt = time.Now().UTC()
		return nil

	} else if attempt.Count == r.maxOTPAttempts {
		return fmt.Errorf("%w: %w", baseErr, ErrReachedMaxLoginAttempts)
	}
	attempt.Count += 1
	return nil
}
func (r *otpAttemptRep) Remove(otpSessionID uuid.UUID) {
	r.attempts.Delete(otpSessionID)
}
