package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/hasher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthAdmin_RegisterAdmin(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	require.Nil(t, err)

	adminCreator := testobj.NewAdminMother()
	hashedPassword := "$2a$10$hashedpassword123"
	password := "password123"
	admin := adminCreator.AdminWithPswdHash(uuid.New(), hashedPassword)
	registerReq := auth.RegisterAdminRequest{
		Adminname: admin.GetUsername(),
		Login:     admin.GetLogin(),
		Password:  password,
	}

	t.Run("success", func(t *testing.T) {
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

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, tokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		err = authAdminServ.RegisterAdmin(ctx, registerReq)
		require.NoError(t, err)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockAdminRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Admin"))
	})

	t.Run("hasher error", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		expectedErr := hasher.ErrHash
		mockHasher.On("HashPassword", password).Return("", expectedErr)

		mockAdminRep := new(adminrep.MockAdminRep)

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, tokenMaker, mockHasher)
		require.Nil(t, err)

		err = authAdminServ.RegisterAdmin(ctx, registerReq)
		require.ErrorIs(t, err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockAdminRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("admin repository error", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockHasher.On("HashPassword", password).Return(hashedPassword, nil)

		mockAdminRep := new(adminrep.MockAdminRep)
		expectedErr := errors.New("database error")
		mockAdminRep.On("Add", ctx, mock.AnythingOfType("*models.Admin")).Return(expectedErr)

		// act
		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, tokenMaker, mockHasher)
		require.Nil(t, err)

		err = authAdminServ.RegisterAdmin(ctx, registerReq)
		require.ErrorIs(t, err, expectedErr)
		mockHasher.AssertCalled(t, "HashPassword", password)
		mockAdminRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Admin"))
	})
}

func TestAuthAdmin_LoginAdmin(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	adminCreator := testobj.NewAdminMother()
	hashedPassword := "$2a$10$hashedpassword123"
	admin := adminCreator.AdminWithPswdHash(uuid.New(), hashedPassword)
	password := "password123"

	loginReq := auth.LoginAdminRequest{
		Login:    admin.GetLogin(),
		Password: password,
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(&admin, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedToken := "access-token-123"
		mockTokenMaker.On("CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration).
			Return(expectedToken, nil)

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)

		require.NoError(t, err)
		require.Equal(t, expectedToken, tokenStr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration)
	})
	t.Run("error admin not found", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		expectedErr := errors.New("admin not found")
		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(nil, expectedErr)

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.ErrorIs(t, err, expectedErr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertNotCalled(t, "CheckPassword", mock.Anything, mock.Anything)
	})

	t.Run("error wrong password", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(&admin, nil)
		expectedErr := errors.New("wrong password")
		mockHasher.On("CheckPassword", password, hashedPassword).Return(expectedErr)

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)
		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.ErrorIs(t, err, expectedErr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
	})

	t.Run("token creation failed", func(t *testing.T) {
		ctx := context.Background()
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		mockAdminRep.On("GetByLogin", ctx, admin.GetLogin()).Return(&admin, nil)
		mockHasher.On("CheckPassword", password, hashedPassword).Return(nil)

		expectedErr := errors.New("token creation failed")
		mockTokenMaker.On("CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration).
			Return("", expectedErr)

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		tokenStr, err := authAdminServ.LoginAdmin(ctx, loginReq)

		require.Error(t, err)
		require.Empty(t, tokenStr)
		require.ErrorIs(t, err, expectedErr)
		mockAdminRep.AssertCalled(t, "GetByLogin", ctx, admin.GetLogin())
		mockHasher.AssertCalled(t, "CheckPassword", password, hashedPassword)
		mockTokenMaker.AssertCalled(t, "CreateToken", admin.GetID(), token.AdminRole, appCnfg.AccessTokenDuration)
	})
}

func TestAuthAdmin_VerifyByToken(t *testing.T) {
	appConfigCreator := testobj.NewAppConfigMother()
	appCnfg := appConfigCreator.Default()

	t.Run("success", func(t *testing.T) {
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

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)

		// act
		payload, err := authAdminServ.VerifyByToken(tokenString)

		require.NoError(t, err)
		require.Equal(t, expectedPayload, payload)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.AdminRole)
	})

	t.Run("error invalid token", func(t *testing.T) {
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "invalid-token"
		expectedErr := errors.New("invalid token")
		mockTokenMaker.On("VerifyToken", tokenString, token.AdminRole).Return(nil, expectedErr)

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)
		// act
		payload, err := authAdminServ.VerifyByToken(tokenString)

		require.Error(t, err)
		require.Nil(t, payload)
		require.ErrorIs(t, err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.AdminRole)
	})

	t.Run("error wrong role", func(t *testing.T) {
		mockHasher := new(hasher.MockHasher)
		mockAdminRep := new(adminrep.MockAdminRep)
		mockTokenMaker := new(token.MockTokenMaker)

		tokenString := "user-token"
		expectedErr := errors.New("invalid role")
		mockTokenMaker.On("VerifyToken", tokenString, token.AdminRole).Return(nil, expectedErr)

		authAdminServ, err := auth.NewAuthAdmin(appCnfg, mockAdminRep, mockTokenMaker, mockHasher)
		require.Nil(t, err)
		// act
		payload, err := authAdminServ.VerifyByToken(tokenString)

		require.Error(t, err)
		require.Nil(t, payload)
		require.ErrorIs(t, err, expectedErr)
		mockTokenMaker.AssertCalled(t, "VerifyToken", tokenString, token.AdminRole)
	})
}
