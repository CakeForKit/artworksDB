package authuser_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	attemptsrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/attempts_rep"
	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
	authsessionrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_session_repository"
	authuser "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_user"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/otp"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/stretchr/testify/mock"
)

type AuthUserServiceSuite struct {
	suite.Suite
}

func TestAuthUserService(t *testing.T) {
	suite.RunSuite(t, new(AuthUserServiceSuite))

}

func (s *AuthUserServiceSuite) TestAuthUser_RegisterUser(t provider.T) {
	durationSession, _ := time.ParseDuration("10m")
	maxAttempts := 5
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	t.Require().NoError(err, "Failed to create token maker")

	userCreator := testobj.NewUserMother()
	hashedPassword := "$2a$10$hashedpassword123"
	user := userCreator.UserWithPswdHash(uuid.New(), hashedPassword)
	passwordUser := "password123"
	registerReq := authmodels.RegisterUserRequest{
		Username:       user.GetUsername(),
		Login:          user.GetLogin(),
		Password:       passwordUser,
		Email:          user.GetEmail(),
		SubscribeEmail: user.IsSubscribedToMail(),
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		mockHasher.On("HashPassword", passwordUser).Return(hashedPassword, nil)

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("Add", ctx, mock.MatchedBy(func(u *models.User) bool {
			return user.GetUsername() == u.GetUsername() &&
				user.GetLogin() == u.GetLogin() &&
				user.GetHashedPassword() == hashedPassword &&
				user.GetEmail() == u.GetEmail() &&
				user.IsSubscribedToMail() == u.IsSubscribedToMail()
		})).Return(nil)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			tokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

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
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			tokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

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
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", passwordUser).Return(hashedPassword, nil)

		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("database error")
		mockUserRep.On("Add", ctx, mock.AnythingOfType("*models.User")).Return(expectedErr)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			tokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

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
	durationSession, _ := time.ParseDuration("10m")
	maxAttempts := 5

	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	userCreator := testobj.NewUserMother()
	hashedPassword := "$2a$10$hashedpassword123"
	user := userCreator.UserWithPswdHash(uuid.New(), hashedPassword)
	passwordUser := "password123"

	loginReq := authmodels.LoginUserRequest{
		Login:    user.GetLogin(),
		Password: passwordUser,
	}

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		// Настройка моков
		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(&user, nil)
		mockHasher.On("CheckPassword", passwordUser, hashedPassword).Return(nil)

		expectedToken := "access-token-123"
		mockTokenMaker.On("CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration).
			Return(expectedToken, nil)

		otpCode := "1234"
		mockOTPServ.On("SendOTP", user).Return(otpCode, nil)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			mockTokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

		// ACT
		sessionID, err := authUserServ.LoginUser(ctx, loginReq)
		sCtx.Require().NoError(err)
		tokenStr, err := authUserServ.OTP(ctx, sessionID, otpCode)
		sCtx.Require().NoError(err)

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
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		expectedErr := errors.New("user not found")
		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(nil, expectedErr)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			mockTokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

		// ACT
		sessionID, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Equal(sessionID, uuid.Nil)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertNotCalled(t, "CheckPassword", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error wrong password", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(&user, nil)

		expectedErr := errors.New("wrong password")
		mockHasher.On("CheckPassword", passwordUser, hashedPassword).Return(expectedErr)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			mockTokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

		// ACT
		sessionID, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Equal(sessionID, uuid.Nil)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", passwordUser, hashedPassword)
	})

	t.WithNewStep("token creation failed", func(sCtx provider.StepCtx) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(&user, nil)
		mockHasher.On("CheckPassword", passwordUser, hashedPassword).Return(nil)

		otpCode := "1234"
		mockOTPServ.On("SendOTP", user).Return(otpCode, nil)

		expectedErr := errors.New("token creation failed")
		mockTokenMaker.On("CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration).
			Return("", expectedErr)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			mockTokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

		// ACT
		sessionID, err := authUserServ.LoginUser(ctx, loginReq)
		sCtx.Require().NoError(err)
		tokenStr, err := authUserServ.OTP(ctx, sessionID, otpCode)

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
	durationSession, _ := time.ParseDuration("10m")
	maxAttempts := 5
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		// ARRANGE
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		expectedPayload := &token.Payload{
			PersonID:  uuid.New(),
			Role:      token.UserRole,
			ExpiredAt: time.Now().Add(time.Hour),
		}

		tokenString := "valid-token-123"
		mockTokenMaker.On("VerifyToken", tokenString, token.UserRole).Return(expectedPayload, nil)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			mockTokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

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
		mockOTPServ := new(otp.MockOTPService)
		authUserSessionRep := authsessionrep.NewAuthUserSessionRep(durationSession)
		loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(maxAttempts, durationSession)
		otpAttemptRep := attemptsrep.NewOTPAttemptRep(maxAttempts, durationSession)

		tokenString := "invalid-token"
		expectedErr := errors.New("invalid token")
		mockTokenMaker.On("VerifyToken", tokenString, token.UserRole).Return(nil, expectedErr)

		authUserServ := authuser.NewAuthUser(
			appCnfg, mockUserRep,
			mockTokenMaker, mockHasher,
			mockOTPServ, authUserSessionRep,
			loginAttemptRep, otpAttemptRep)

		// ACT
		payload, err := authUserServ.VerifyByToken(tokenString)

		// ASSERT
		sCtx.Require().Error(err)
		sCtx.Assert().Nil(payload)
		sCtx.Assert().Equal(expectedErr, err)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.UserRole)
	})
}
