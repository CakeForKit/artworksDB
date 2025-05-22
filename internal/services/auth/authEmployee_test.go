package auth

import (
	"context"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestConfig() cnfg.AppConfig {
	config := cnfg.AppConfig{
		TokenSymmetricKey:   "01234567890123456789012345678912",
		AccessTokenDuration: time.Hour,
	}
	return config
}

func createTestEmployee() *models.Employee {
	employee, _ := models.NewEmployee(
		uuid.New(),
		"test_user",
		"test_login",
		"hashed_password",
		time.Now(),
		true,
		uuid.New(),
	)
	return &employee
}

func TestAuthEmployee_LoginEmployee(t *testing.T) {
	ctx := context.Background()
	config := createTestConfig()
	validPassword := "valid_password"
	invalidPassword := "invalid_password"

	tests := []struct {
		name          string
		login         string
		password      string
		mockEmployee  *models.Employee
		mockError     error
		checkPassword error
		expectedError error
	}{
		{
			name:          "success",
			login:         "test_login",
			password:      validPassword,
			mockEmployee:  createTestEmployee(),
			mockError:     nil,
			checkPassword: nil,
			expectedError: nil,
		},
		{
			name:          "employee not found",
			login:         "unknown_login",
			password:      validPassword,
			mockEmployee:  nil,
			mockError:     employeerep.ErrEmployeeNotFound,
			expectedError: employeerep.ErrEmployeeNotFound,
		},
		{
			name:          "invalid password",
			login:         "test_login",
			password:      invalidPassword,
			mockEmployee:  createTestEmployee(),
			mockError:     nil,
			checkPassword: hasher.ErrPassword,
			expectedError: hasher.ErrPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(employeerep.MockEmployeeRep)
			tokenMaker, _ := token.NewTokenMaker(config.TokenSymmetricKey)
			hasher := new(MockHasher)

			service := &authEmployee{
				tokenMaker:  tokenMaker,
				config:      config,
				employeerep: mockRepo,
				hasher:      hasher,
			}

			mockRepo.On("GetByLogin", ctx, tt.login).Return(tt.mockEmployee, tt.mockError)
			if tt.mockEmployee != nil {
				hasher.On("CheckPassword", tt.password, tt.mockEmployee.GetHashedPassword()).Return(tt.checkPassword)
			}

			token, err := service.LoginEmployee(ctx, LoginEmployeeRequest{
				Login:    tt.login,
				Password: tt.password,
			})

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
		})
	}
}

func TestAuthEmployee_RegisterEmployee(t *testing.T) {
	ctx := context.Background()
	config := createTestConfig()
	validRequest := RegisterEmployeeRequest{
		Username: "new_user",
		Login:    "new_login",
		Password: "new_password",
	}
	adminID := uuid.New()

	tests := []struct {
		name          string
		request       RegisterEmployeeRequest
		hashError     error
		addError      error
		expectedError error
	}{
		{
			name:          "success",
			request:       validRequest,
			hashError:     nil,
			addError:      nil,
			expectedError: nil,
		},
		{
			name:          "hash error",
			request:       validRequest,
			hashError:     hasher.ErrHash,
			addError:      nil,
			expectedError: hasher.ErrHash,
		},
		{
			name:          "add error",
			request:       validRequest,
			hashError:     nil,
			addError:      employeerep.ErrFailedToAddEmployee,
			expectedError: employeerep.ErrFailedToAddEmployee,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(employeerep.MockEmployeeRep)
			tokenMaker, _ := token.NewTokenMaker(config.TokenSymmetricKey)
			hasher := new(MockHasher)

			hasher.On("HashPassword", tt.request.Password).Return("hashed_password", tt.hashError)
			if tt.hashError == nil {
				mockRepo.On("Add", ctx, mock.Anything).Return(tt.addError)
			}

			service := &authEmployee{
				tokenMaker:  tokenMaker,
				config:      config,
				employeerep: mockRepo,
				hasher:      hasher,
			}

			err := service.RegisterEmployee(ctx, tt.request, adminID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
		})
	}
}
