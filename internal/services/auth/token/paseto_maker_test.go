package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	userID := uuid.New()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(userID, UserRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token, UserRole)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.GetID())
	require.Equal(t, userID, payload.GetPersonID())
	require.WithinDuration(t, issuedAt, payload.GetIssuedAt(), time.Second)
	require.WithinDuration(t, expiredAt, payload.GetExpiredAt(), time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	token, err := maker.CreateToken(uuid.New(), UserRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token, UserRole)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
