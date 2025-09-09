package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthEmployee_RegisterEmployee(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	require.Nil(t, err)

	employeeCreator := testobj.NewEmployeeMother()
	hashedPassword := "$2a$10$hashedpassword123"
	password := "password123"
	adminID := uuid.New()
	employee := employeeCreator.EmployeeWithPswdHash(uuid.New(), adminID, hashedPassword)
	registerReq := auth.RegisterEmployeeRequest{
		Username: employee.GetUsername(),
		Login:    employee.GetLogin(),
		Password: password,
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", password).Return(hashedPassword, nil)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockEmployeeRep.On("Add", ctx, mock.MatchedBy(func(u *models.Employee) bool {
			return employee.GetUsername() == u.GetUsername() &&
				employee.GetLogin() == u.GetLogin() &&
				employee.GetHashedPassword() == hashedPassword &&
				employee.GetAdminID() == u.GetAdminID() &&
				employee.IsValid() == u.IsValid()
		})).Return(nil)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, tokenMaker, mockHasher)
		require.Nil(t, err)
		_ = authEmployee
		_ = registerReq

		// act
		err = authEmployee.RegisterEmployee(ctx, registerReq, adminID)
		require.NoError(t, err)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockEmployeeRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Employee"))
	})

	t.Run("hasher error", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		expectedErr := hasher.ErrHash
		mockHasher.On("HashPassword", password).Return("", expectedErr)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, tokenMaker, mockHasher)
		require.Nil(t, err)

		err = authEmployee.RegisterEmployee(ctx, registerReq, adminID)
		require.ErrorIs(t, err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockEmployeeRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})
	t.Run("employee repository error", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", password).Return(hashedPassword, nil)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		expectedErr := errors.New("database error")
		mockEmployeeRep.On("Add", ctx, mock.AnythingOfType("*models.Employee")).Return(expectedErr)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, tokenMaker, mockHasher)
		require.Nil(t, err)

		err = authEmployee.RegisterEmployee(ctx, registerReq, adminID)
		require.ErrorIs(t, err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockEmployeeRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Employee"))
	})
}

func TestAuthEmployee_LoginEmployee(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	employeeCreator := testobj.NewEmployeeMother()
	hashedPassword := "$2a$10$hashedpassword123"
	adminID := uuid.New()
	employee := employeeCreator.EmployeeWithPswdHash(uuid.New(), adminID, hashedPassword)
	password := "password123"

	loginReq := auth.LoginEmployeeRequest{
		Login:    employee.GetLogin(),
		Password: password,
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&employee, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedToken := "access-token-123"
		mockTokenMaker.On("CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration).
			Return(expectedToken, nil)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)

		require.NoError(t, err)
		require.Equal(t, expectedToken, tokenStr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration)
	})
	t.Run("error employee not found", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedErr := errors.New("employee not found")
		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(nil, expectedErr)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.ErrorIs(t, err, expectedErr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
	})

	t.Run("error employee not valid", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		invalidEmployee := employeeCreator.EmployeeWithPswdHash(uuid.New(), adminID, hashedPassword)
		invalidEmployee.SetValid(false)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&invalidEmployee, nil)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.ErrorIs(t, err, auth.ErrEmployeeNotValid)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertNotCalled(t, "CheckPassword", mock.Anything, mock.Anything)
	})

	t.Run("error wrong password", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&employee, nil)
		expectedErr := errors.New("wrong password")
		mockHasher.On("CheckPassword", password, hashedPassword).Return(expectedErr)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.ErrorIs(t, err, expectedErr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
	})

	t.Run("token creation failed", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&employee, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedErr := errors.New("token creation failed")
		mockTokenMaker.On("CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration).
			Return("", expectedErr)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)

		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.ErrorIs(t, err, expectedErr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration)
	})
}

func TestAuthEmployee_VerifyByToken(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	t.Run("success", func(t *testing.T) {
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedPayload := &token.Payload{
			PersonID:  uuid.New(),
			Role:      token.EmployeeRole,
			ExpiredAt: time.Now().Add(time.Hour),
		}

		tokenString := "valid-token-123"
		mockTokenMaker.On("VerifyToken", tokenString, token.EmployeeRole).Return(expectedPayload, nil)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		payload, err := authEmployee.VerifyByToken(tokenString)

		require.NoError(t, err)
		require.Equal(t, expectedPayload, payload)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.EmployeeRole)
	})

	t.Run("error invalid token", func(t *testing.T) {
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "invalid-token"
		expectedErr := errors.New("invalid token")
		mockTokenMaker.On("VerifyToken", tokenString, token.EmployeeRole).Return(nil, expectedErr)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		payload, err := authEmployee.VerifyByToken(tokenString)

		require.Error(t, err)
		require.Nil(t, payload)
		require.ErrorIs(t, err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.EmployeeRole)
	})

	t.Run("error wrong role", func(t *testing.T) {
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "admin-token"
		expectedErr := errors.New("invalid role")
		mockTokenMaker.On("VerifyToken", tokenString, token.EmployeeRole).Return(nil, expectedErr)

		authEmployee, err := auth.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		payload, err := authEmployee.VerifyByToken(tokenString)

		require.Error(t, err)
		require.Nil(t, payload)
		require.ErrorIs(t, err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.EmployeeRole)
	})

}
