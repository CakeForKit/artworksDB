package auth

import (
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep/mockemployeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/util"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
	"github.com/stateio/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestAuthEmployee(t *testing.T) {
	validEmployeeID := uuid.New()
	validUsername := "employee_user"
	validLogin := "employee_login"
	validPassword := "employee_password"
	hashedPassword := "hashed_employee_password"

	config := util.Config{
		TokenSymmetricKey:    "01234567890123456789012345678912",
		AccessTokenDuration:  time.Hour,
		RefreshTokenDuration: time.Hour * 24,
	}

	t.Run("LoginEmployee", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			employeeRep := new(mockemployeerep.MockEmployeeRep)
			tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			hasher := new(MockHasher)

			employee, err := models.NewEmployee(
				validEmployeeID,
				validUsername,
				validLogin,
				hashedPassword,
				time.Now(),
			)
			require.NoError(t, err)

			employeeRep.On("GetByLogin", validLogin).Return(&employee, nil)
			hasher.On("CheckPassword", validPassword, hashedPassword).Return(nil)

			service := &authEmployee{
				tokenMaker:  tokenMaker,
				config:      config,
				employeerep: employeeRep,
				hasher:      hasher,
			}

			_, err = service.LoginEmployee(LoginEmployeeRequest{
				Login:    validLogin,
				Password: validPassword,
			})
			require.NoError(t, err)

			employeeRep.AssertExpectations(t)
			hasher.AssertExpectations(t)
		})

		t.Run("EmployeeNotFound", func(t *testing.T) {
			employeeRep := new(mockemployeerep.MockEmployeeRep)
			tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			hasher := new(MockHasher)

			var nile *models.Employee = nil
			employeeRep.On("GetByLogin", validLogin).Return(nile, assert.AnError)

			service := &authEmployee{
				tokenMaker:  tokenMaker,
				config:      config,
				employeerep: employeeRep,
				hasher:      hasher,
			}

			_, err = service.LoginEmployee(LoginEmployeeRequest{
				Login:    validLogin,
				Password: validPassword,
			})

			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("InvalidPassword", func(t *testing.T) {
			employeeRep := new(mockemployeerep.MockEmployeeRep)
			tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			hasher := new(MockHasher)

			employee, err := models.NewEmployee(
				validEmployeeID,
				validUsername,
				validLogin,
				hashedPassword,
				time.Now(),
			)
			require.NoError(t, err)

			employeeRep.On("GetByLogin", validLogin).Return(&employee, nil)
			hasher.On("CheckPassword", validPassword, hashedPassword).Return(assert.AnError)

			service := &authEmployee{
				tokenMaker:  tokenMaker,
				config:      config,
				employeerep: employeeRep,
				hasher:      hasher,
			}

			_, err = service.LoginEmployee(LoginEmployeeRequest{
				Login:    validLogin,
				Password: validPassword,
			})

			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("RegisterEmployee", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			employeeRep := new(mockemployeerep.MockEmployeeRep)
			tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			hasher := new(MockHasher)

			hasher.On("HashPassword", validPassword).Return(hashedPassword, nil)
			employeeRep.On("Add", mock.AnythingOfType("*models.Employee")).Return(nil)

			service := &authEmployee{
				tokenMaker:  tokenMaker,
				config:      config,
				employeerep: employeeRep,
				hasher:      hasher,
			}

			err = service.RegisterEmployee(RegisterEmployeeRequest{
				Username: validUsername,
				Login:    validLogin,
				Password: validPassword,
			})

			require.NoError(t, err)
			hasher.AssertExpectations(t)
			employeeRep.AssertExpectations(t)
		})
	})
}
