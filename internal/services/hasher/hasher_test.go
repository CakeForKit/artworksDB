package hasher_test

import (
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/hasher"
	"github.com/stateio/testify/assert"
)

func TestPassword(t *testing.T) {
	password := "1234were### _"
	hasherServ, err := hasher.NewHasher()
	assert.NoError(t, err)

	hashedPassword, err := hasherServ.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	err = hasherServ.CheckPassword(password, hashedPassword)
	assert.NoError(t, err)
}

func TestWrongPassword(t *testing.T) {
	password := "1234were### _"
	hasherServ, err := hasher.NewHasher()
	assert.NoError(t, err)

	hashedPassword, err := hasherServ.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	wrongPassword := "0"
	err = hasherServ.CheckPassword(wrongPassword, hashedPassword)
	assert.Error(t, err)
	assert.Equal(t, hasher.ErrPassword, err)
}

func TestDiffHash(t *testing.T) {
	password := "1234"
	hasherServ, err := hasher.NewHasher()
	assert.NoError(t, err)

	hashedPassword1, err := hasherServ.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword1)

	hashedPassword2, err := hasherServ.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword2)

	assert.NotEqual(t, hashedPassword1, hashedPassword2)
}
