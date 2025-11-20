package authadmin_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	authadmin "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_admin"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type AuthAdminServiceSuite struct {
	suite.Suite
}

func TestAuthAdminService(t *testing.T) {
	suite.RunSuite(t, new(AuthAdminServiceSuite))
}

func (s *AuthAdminServiceSuite) TestAuthAdmin_RegisterAdmin(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	t.Require().NoError(err, "Failed to create token maker")

	adminCreator := testobj.NewAdminMother()
	hashedPassword := "$2a$10$hashedpassword123"
	password := "password123"
	admin := adminCreator.AdminWithPswdHash(uuid.New(), hashedPassword)
	registerReq := authadmin.RegisterAdminRequest{
		Adminname: admin.GetUsername(),
		Login:     admin.GetLogin(),
		Password:  password,
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", password).Return(hashedPassword, nil)

		mockAdminRep := new(adminrep.MockAdminRep)
		mockAdminRep.On("Add", ctx, mock.MatchedBy(func(a *models.Admin) bool {
			return admin.GetUsername() == a.GetUsername() &&
				admin.GetLogin() == a.GetLogin() &&
				admin.GetHashedPassword() == hashedPassword &&
				admin.IsValid() == a.IsValid()
		})).Return(nil)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, tokenMaker, mockHasher)

		// act
		err = authAdminServ.RegisterAdmin(ctx, registerReq)
		sCtx.Require().NoError(err)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockAdminRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Admin"))
	})

	t.WithNewStep("hasher error", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		expectedErr := hasher.ErrHash
		mockHasher.On("HashPassword", password).Return("", expectedErr)

		mockAdminRep := new(adminrep.MockAdminRep)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, tokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		err = authAdminServ.RegisterAdmin(ctx, registerReq)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockAdminRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.WithNewStep("admin repository error", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", password).Return(hashedPassword, nil)

		mockAdminRep := new(adminrep.MockAdminRep)
		expectedErr := errors.New("database error")
		mockAdminRep.On("Add", ctx, mock.AnythingOfType("*models.Admin")).Return(expectedErr)

		// act
		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, tokenMaker, mockHasher)

		err = authAdminServ.RegisterAdmin(ctx, registerReq)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockAdminRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Admin"))
	})
}

func (s *AuthAdminServiceSuite) TestAuthAdmin_LoginAdmin(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	adminCreator := testobj.NewAdminMother()
	hashedPassword := "$2a$10$hashedpassword123"
	admin := adminCreator.AdminWithPswdHash(uuid.New(), hashedPassword)
	password := "password123"

	loginReq := authadmin.LoginAdminRequest{
		Login:    admin.GetLogin(),
		Password: password,
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(&admin, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedToken := "access-token-123"
		mockTokenMaker.On("CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration).
			Return(expectedToken, nil)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)

		// act
		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)

		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(expectedToken, tokenStr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration)
	})
	t.WithNewStep("error admin not found", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedErr := errors.New("admin not found")
		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(nil, expectedErr)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)

		// act
		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertNotCalled(t, "CheckPassword", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error wrong password", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(&admin, nil)
		expectedErr := errors.New("wrong password")
		mockHasher.On("CheckPassword", password, hashedPassword).Return(expectedErr)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)

		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
	})

	t.WithNewStep("token creation failed", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(&admin, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedErr := errors.New("token creation failed")
		mockTokenMaker.On("CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration).
			Return("", expectedErr)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)

		// act
		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)

		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration)
	})
}

func (s *AuthAdminServiceSuite) TestAuthAdmin_VerifyByToken(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedPayload := &token.Payload{
			PersonID:  uuid.New(),
			Role:      token.AdminRole,
			ExpiredAt: time.Now().Add(time.Hour),
		}

		tokenString := "valid-token-123"
		mockTokenMaker.On("VerifyToken", tokenString, token.AdminRole).Return(expectedPayload, nil)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)

		// act
		payload, err := authAdminServ.VerifyByToken(tokenString)

		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(expectedPayload, payload)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.AdminRole)
	})

	t.WithNewStep("error invalid token", func(sCtx provider.StepCtx) {
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "invalid-token"
		expectedErr := errors.New("invalid token")
		mockTokenMaker.On("VerifyToken", tokenString, token.AdminRole).Return(nil, expectedErr)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)

		// act
		payload, err := authAdminServ.VerifyByToken(tokenString)

		sCtx.Require().Error(err)
		sCtx.Assert().Nil(payload)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.AdminRole)
	})

	t.WithNewStep("error wrong role", func(sCtx provider.StepCtx) {
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "user-token"
		expectedErr := errors.New("invalid role")
		mockTokenMaker.On("VerifyToken", tokenString, token.AdminRole).Return(nil, expectedErr)

		authAdminServ := authadmin.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)

		// act
		payload, err := authAdminServ.VerifyByToken(tokenString)

		sCtx.Require().Error(err)
		sCtx.Assert().Nil(payload)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.AdminRole)
	})
}
