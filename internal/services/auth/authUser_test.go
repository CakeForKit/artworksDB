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

	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthUser_RegisterUser(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	require.Nil(t, err)

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

	t.Run("success", func(t *testing.T) {
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
		require.Nil(t, err)
		// act
		err = authUserServ.RegisterUser(ctx, registerReq)

		require.NoError(t, err)
		mockHasher.AssertCalled(t, "HashPassword", passwordUser)
		mockUserRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.User"))
	})
	t.Run("hasher error", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		expectedErr := hasher.ErrHash
		mockHasher.On("HashPassword", passwordUser).Return("", expectedErr)

		mockUserRep := new(userrep.MockUserRep)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, tokenMaker, mockHasher)
		require.Nil(t, err)

		// ACT
		err = authUserServ.RegisterUser(ctx, registerReq)

		// ASSERT
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		mockHasher.AssertCalled(t, "HashPassword", passwordUser)
		mockUserRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("user repository error", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", passwordUser).Return(hashedPassword, nil)

		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("database error")
		mockUserRep.On("Add", ctx, mock.AnythingOfType("*models.User")).Return(expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, tokenMaker, mockHasher)
		require.Nil(t, err)

		// ACT
		err = authUserServ.RegisterUser(ctx, registerReq)

		// ASSERT
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		mockHasher.AssertCalled(t, "HashPassword", passwordUser)
		mockUserRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.User"))
	})
}

func TestAuthUser_LoginUser(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
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
		require.Nil(t, err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expectedToken, tokenStr)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", passwordUser, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration)
	})

	t.Run("error user not found", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedErr := errors.New("user not found")
		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(nil, expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.Equal(t, expectedErr, err)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertNotCalled(t, "CheckPassword", mock.Anything, mock.Anything)
	})

	t.Run("error wrong password", func(t *testing.T) {
		// ARRANGE
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockUserRep.On("GetByLogin", ctx, user.GetLogin()).Return(&user, nil)

		expectedErr := errors.New("wrong password")
		mockHasher.On("CheckPassword", passwordUser, hashedPassword).Return(expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.Equal(t, expectedErr, err)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", passwordUser, hashedPassword)
	})

	t.Run("token creation failed", func(t *testing.T) {
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
		require.Nil(t, err)

		// ACT
		tokenStr, err := authUserServ.LoginUser(ctx, loginReq)

		// ASSERT
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.Equal(t, expectedErr, err)
		mockUserRep.AssertCalled(t, "GetByLogin", ctx, user.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", passwordUser, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", user.GetID(), token.UserRole, appCnfg.AccessTokenDuration)
	})
}

func TestAuthUser_VerifyByToken(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	t.Run("success", func(t *testing.T) {
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
		require.Nil(t, err)

		// ACT
		payload, err := authUserServ.VerifyByToken(tokenString)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expectedPayload, payload)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.UserRole)
	})

	t.Run("error", func(t *testing.T) {
		// ARRANGE
		mockHasher := new(hasher.MockHasher)
		mockUserRep := new(userrep.MockUserRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "invalid-token"
		expectedErr := errors.New("invalid token")
		mockTokenMaker.On("VerifyToken", tokenString, token.UserRole).Return(nil, expectedErr)

		authUserServ, err := auth.NewAuthUser(appCnfg, mockUserRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// ACT
		payload, err := authUserServ.VerifyByToken(tokenString)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, payload)
		require.Equal(t, expectedErr, err)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.UserRole)
	})
}
