package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/mock"
)

type AuthUserServiceSuite struct {
	suite.Suite
}

func TestAuthUserService(t *testing.T) {
	suite.RunSuite(t, new(AuthUserServiceSuite))

}

func (s *AuthUserServiceSuite) TestAuthUser_RegisterUser(t provider.T) {

	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	t.Require().NoError(err, "Failed to create token maker")

	userCreator := testobj.NewUserMother()
	hashedPassword := "$2a$10$hashedpassword123"
	user := userCreator.UserWithPswdHash(uuid.New(), hashedPassword)
	passwordUser := "password123"
	registerReq := auth.RegisterUserRequest{
		Username:       user.GetUsername(),
		Login:          user.GetLogin(),
		Password:       passwordUser,
		Email:          user.GetEmail(),
		SubscribeEmail: user.IsSubscribedToMail(),
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)

		mockHasher.On("HashPassword", passwordUser).Return(hashedPassword, nil)

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("Add", ctx, mock.MatchedBy(func(u *models.User) bool {
			return user.GetUsername() == u.GetUsername() &&
				user.GetLogin() == u.GetLogin() &&
				user.GetHashedPassword() == hashedPassword &&
				user.GetEmail() == u.GetEmail() &&
				user.IsSubscribedToMail() == u.IsSubscribedToMail()
		})).Return(nil)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, tokenMaker, mockHasher)
		sCtx.Require().NoError(err)
		// act
		err = authUserServ.RegisterUser(ctx, registerReq)

		sCtx.Require().NoError(err)
		mockHasher.AssertCalled(t, "HashPassword", passwordUser)
		mockUserRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.User"))
	})
	t.WithNewStep("hasher error", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		expectedErr := hasher.ErrHash
		mockHasher.On("HashPassword", passwordUser).Return("", expectedErr)

		mockUserRep := new(userrep.MockUserRep)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, tokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		err = authUserServ.RegisterUser(ctx, registerReq)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Equal(expectedErr, err)
		mockHasher.AssertCalled(t, "HashPassword", passwordUser)
		mockUserRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.WithNewStep("user repository error", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", passwordUser).Return(hashedPassword, nil)

		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("database error")
		mockUserRep.On("Add", ctx, mock.AnythingOfType("*models.User")).Return(expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, tokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		err = authUserServ.RegisterUser(ctx, registerReq)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Equal(expectedErr, err)
		mockHasher.AssertCalled(t, "HashPassword", passwordUser)
		mockUserRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.User"))
	})
}

func (s *AuthUserServiceSuite) TestAuthUser_LoginUser(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	userCreator := testobj.NewUserMother()
	hashedPassword := "$2a$10$hashedpassword123"
	user := userCreator.UserWithPswdHash(uuid.New(), hashedPassword)
	passwordUser := "password123"

	loginReq := auth.LoginUserRequest{
		Login:    user.GetLogin(),
		Password: passwordUser,
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		// Настройка моков
		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(&user, nil)
		mockHasher.On("CheckPassword", passwordUser, hashedPassword).Return(nil)

		expectedToken := "access-token-123"
		mockTokenMaker.On("CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration).
			Return(expectedToken, nil)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(expectedToken, tokenStr)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", passwordUser, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration)
	})

	t.WithNewStep("error user not found", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedErr := errors.New("user not found")
		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(nil, expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Assert().Equal(expectedErr, err)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertNotCalled(t, "CheckPassword", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error wrong password", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(&user, nil)

		expectedErr := errors.New("wrong password")
		mockHasher.On("CheckPassword", passwordUser, hashedPassword).Return(expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Assert().Equal(expectedErr, err)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", passwordUser, hashedPassword)
	})

	t.WithNewStep("token creation failed", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(&user, nil)
		mockHasher.On("CheckPassword", passwordUser, hashedPassword).Return(nil)

		expectedErr := errors.New("token creation failed")
		mockTokenMaker.On("CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration).
			Return("", expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Empty(tokenStr)
		sCtx.Assert().Equal(expectedErr, err)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", passwordUser, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration)
	})
}

func (s *AuthUserServiceSuite) TestAuthUser_VerifyByToken(t provider.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		// ARRANGE
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedPayload := &token.Payload{
			PersonID:  uuid.New(),
			Role:      token.UserRole,
			ExpiredAt: time.Now().Add(time.Hour),
		}

		tokenString := "valid-token-123"
		mockTokenMaker.On("VerifyToken", tokenString, token.UserRole).Return(expectedPayload, nil)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		payload, err := authUserServ.VerifyByToken(tokenString)

		// ASSERT
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(expectedPayload, payload)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.UserRole)
	})

	t.WithNewStep("error", func(sCtx provider.StepCtx) {
		// ARRANGE
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "invalid-token"
		expectedErr := errors.New("invalid token")
		mockTokenMaker.On("VerifyToken", tokenString, token.UserRole).Return(nil, expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		sCtx.Require().NoError(err)

		// ACT
		payload, err := authUserServ.VerifyByToken(tokenString)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Nil(payload)
		sCtx.Assert().Equal(expectedErr, err)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.UserRole)
	})
}
