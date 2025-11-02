package attemptsrep

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type LoginAttemptUserRep interface {
	Add(login string) error
	Remove(login string)
}

type loginAttempt struct {
	LoginKey  string
	Count     int
	CreatedAt time.Time
}

var (
	ErrLoginAttemptUserRep     = errors.New("LoginAttemptUserRep")
	ErrSerialyzeSession        = errors.New("error serialyze session value")
	ErrReachedMaxLoginAttempts = errors.New("reached max attempts")
)

type loginAttemptUserRep struct {
	maxLoginAttempts int
	durationSession  time.Duration
	attempts         sync.Map
}

func NewLoginAttemptUserRep(maxLoginAttempts int, durationSession time.Duration) LoginAttemptUserRep {
	return &loginAttemptUserRep{
		maxLoginAttempts: maxLoginAttempts,
		durationSession:  durationSession,
		attempts:         sync.Map{},
	}
}

func (r *loginAttemptUserRep) Add(login string) error {
	baseErr := fmt.Errorf("%w: Add", ErrLoginAttemptUserRep)
	val, ok := r.attempts.Load(login)
	if !ok {
		r.attempts.Store(login, &loginAttempt{
			LoginKey:  login,
			Count:     1,
			CreatedAt: time.Now().UTC(),
		})
		return nil
	}
	attempt, ok := val.(*loginAttempt)
	if !ok {
		return fmt.Errorf("%w: %w", baseErr, ErrSerialyzeSession)
	}
	if attempt.CreatedAt.Sub(time.Now().UTC()) > r.durationSession {
		attempt.Count = 1
		attempt.CreatedAt = time.Now().UTC()
		return nil

	} else if attempt.Count == r.maxLoginAttempts {
		return fmt.Errorf("%w: %w", baseErr, ErrReachedMaxLoginAttempts)
	}
	attempt.Count += 1
	return nil
}

func (r *loginAttemptUserRep) Remove(login string) {
	r.attempts.Delete(login)
}
