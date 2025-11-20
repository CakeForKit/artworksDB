package authemployee_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	authemployee "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_employee"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type AuthEmployeeServiceSuite struct {
	suite.Suite
}

func TestAuthEmployeeService(t *testing.T) {
	suite.RunSuite(t, new(AuthEmployeeServiceSuite))
}

func (s *AuthEmployeeServiceSuite) TestAuthEmployee_RegisterEmployee(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	t.Require().NoError(err, "Failed to create token maker")

	employeeCreator := testobj.NewEmployeeMother()
	hashedPassword := "$2a$10$hashedpassword123"
	password := "password123"
	adminID := uuid.New()
	employee := employeeCreator.EmployeeWithPswdHash(uuid.New(), adminID, hashedPassword)
	registerReq := authemployee.RegisterEmployeeRequest{
		Username: employee.GetUsername(),
		Login:    employee.GetLogin(),
		Password: password,
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
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

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, tokenMaker, mockHasher)

		_ = authEmployee
		_ = registerReq

		// act
		err = authEmployee.RegisterEmployee(ctx, registerReq, adminID)
		sCtx.Require().NoError(err)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockEmployeeRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Employee"))
	})

	t.WithNewStep("hasher error", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		expectedErr := hasher.ErrHash
		mockHasher.On("HashPassword", password).Return("", expectedErr)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, tokenMaker, mockHasher)

		err = authEmployee.RegisterEmployee(ctx, registerReq, adminID)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockEmployeeRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})
	t.WithNewStep("employee repository error", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", password).Return(hashedPassword, nil)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		expectedErr := errors.New("database error")
		mockEmployeeRep.On("Add", ctx, mock.AnythingOfType("*models.Employee")).Return(expectedErr)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, tokenMaker, mockHasher)

		err = authEmployee.RegisterEmployee(ctx, registerReq, adminID)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockEmployeeRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Employee"))
	})
}

func (s *AuthEmployeeServiceSuite) TestAuthEmployee_LoginEmployee(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	employeeCreator := testobj.NewEmployeeMother()
	hashedPassword := "$2a$10$hashedpassword123"
	adminID := uuid.New()
	employee := employeeCreator.EmployeeWithPswdHash(uuid.New(), adminID, hashedPassword)
	password := "password123"

	loginReq := authemployee.LoginEmployeeRequest{
		Login:    employee.GetLogin(),
		Password: password,
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&employee, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedToken := "access-token-123"
		mockTokenMaker.On("CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration).
			Return(expectedToken, nil)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)

		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(expectedToken, tokenStr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration)
	})
	t.WithNewStep("error employee not found", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedErr := errors.New("employee not found")
		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(nil, expectedErr)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
	})

	t.WithNewStep("error employee not valid", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		invalidEmployee := employeeCreator.EmployeeWithPswdHash(uuid.New(), adminID, hashedPassword)
		invalidEmployee.SetValid(false)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&invalidEmployee, nil)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Require().ErrorIs(err, authemployee.ErrEmployeeNotValid)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertNotCalled(t, "CheckPassword", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error wrong password", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&employee, nil)
		expectedErr := errors.New("wrong password")
		mockHasher.On("CheckPassword", password, hashedPassword).Return(expectedErr)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
	})

	t.WithNewStep("token creation failed", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockEmployeeRep.On("GetByLogin", ctx, employee.GetLogin()).Return(&employee, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedErr := errors.New("token creation failed")
		mockTokenMaker.On("CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration).
			Return("", expectedErr)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		// act
		tokenStr, err := authEmployee.LoginEmployee(ctx, loginReq)

		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockEmployeeRep.AssertCalled(t, "GetByLogin", ctx, employee.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", employee.GetID(), token.EmployeeRole, appCnfg.AccessTokenDuration)
	})
}

func (s *AuthEmployeeServiceSuite) TestAuthEmployee_VerifyByToken(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
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

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		// act
		payload, err := authEmployee.VerifyByToken(tokenString)

		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(expectedPayload, payload)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.EmployeeRole)
	})

	t.WithNewStep("error invalid token", func(sCtx provider.StepCtx) {
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "invalid-token"
		expectedErr := errors.New("invalid token")
		mockTokenMaker.On("VerifyToken", tokenString, token.EmployeeRole).Return(nil, expectedErr)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		// act
		payload, err := authEmployee.VerifyByToken(tokenString)

		sCtx.Require().Error(err)
		sCtx.Assert().Nil(payload)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.EmployeeRole)
	})

	t.WithNewStep("error wrong role", func(sCtx provider.StepCtx) {
		mockHasher := new(hasher.MockHasher)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "admin-token"
		expectedErr := errors.New("invalid role")
		mockTokenMaker.On("VerifyToken", tokenString, token.EmployeeRole).Return(nil, expectedErr)

		authEmployee := authemployee.NewAuthEmployee(appCnfg, mockEmployeeRep, mockTokenMaker, mockHasher)

		// act
		payload, err := authEmployee.VerifyByToken(tokenString)

		sCtx.Require().Error(err)
		sCtx.Assert().Nil(payload)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.EmployeeRole)
	})

}
