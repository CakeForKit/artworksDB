package auth

import (
	"context"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
	"github.com/stateio/testify/require"
	"github.com/stretchr/testify/assert"
)

// MockHasher implements hasher.Hasher interface
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) CheckPassword(password string, hashedPassword string) error {
	args := m.Called(password, hashedPassword)
	return args.Error(0)
}

func TestAuthUserService(t *testing.T) {
	validUserID := uuid.New()
	validUsername := "testuser"
	validLogin := "testlogin"
	validPassword := "securepassword"
	validEmail := "test@example.com"
	hashedPassword := "hashedpassword"

	config := cnfg.AppConfig{
		TokenSymmetricKey:   "01234567890123456789012345678912",
		AccessTokenDuration: time.Hour,
	}

	t.Run("LoginUser", func(t *testing.T) {
		ctx := context.Background()
		t.Run("Success", func(t *testing.T) {
			userRep := new(userrep.MockUserRep)
			tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			hasher := new(MockHasher)

			user, err := models.NewUser(
				validUserID,
				validUsername,
				validLogin,
				hashedPassword,
				time.Now(),
				validEmail,
				true,
			)
			require.NoError(t, err)

			userRep.On("GetByLogin", ctx, validLogin).Return(&user, nil)
			hasher.On("CheckPassword", validPassword, hashedPassword).Return(nil)

			service := &authUser{
				tokenMaker: tokenMaker,
				config:     config,
				userrep:    userRep,
				hasher:     hasher,
			}

			_, err = service.LoginUser(ctx, LoginUserRequest{
				Login:    validLogin,
				Password: validPassword,
			})
			require.NoError(t, err)

			userRep.AssertExpectations(t)
			hasher.AssertExpectations(t)
		})

		t.Run("InvalidPassword", func(t *testing.T) {
			userRep := new(userrep.MockUserRep)
			tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			hasher := new(MockHasher)

			user, err := models.NewUser(
				validUserID,
				validUsername,
				validLogin,
				hashedPassword,
				time.Now(),
				validEmail,
				true,
			)
			require.NoError(t, err)

			userRep.On("GetByLogin", ctx, validLogin).Return(&user, nil)
			hasher.On("CheckPassword", validPassword, hashedPassword).Return(assert.AnError)

			service := &authUser{
				tokenMaker: tokenMaker,
				config:     config,
				userrep:    userRep,
				hasher:     hasher,
			}

			_, err = service.LoginUser(ctx, LoginUserRequest{
				Login:    validLogin,
				Password: validPassword,
			})

			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

	})

	t.Run("RegisterUser", func(t *testing.T) {
		ctx := context.Background()
		t.Run("Success", func(t *testing.T) {
			userRep := new(userrep.MockUserRep)
			tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			hasher := new(MockHasher)

			hasher.On("HashPassword", validPassword).Return(hashedPassword, nil)
			userRep.On("Add", ctx, mock.Anything).Return(nil)

			service := &authUser{
				tokenMaker: tokenMaker,
				config:     config,
				userrep:    userRep,
				hasher:     hasher,
			}

			err = service.RegisterUser(ctx, RegisterUserRequest{
				Username:       validUsername,
				Login:          validLogin,
				Password:       validPassword,
				Email:          validEmail,
				SubscribeEmail: true,
			})

			require.NoError(t, err)

			hasher.AssertExpectations(t)
			userRep.AssertExpectations(t)
		})
	})
}
