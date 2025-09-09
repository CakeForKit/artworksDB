package userservice_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/userservice"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetSelf(t *testing.T) {
	userCreator := testobj.NewUserMother()
	user := userCreator.DefaultUser(uuid.New())

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetByID", ctx, user.GetID()).Return(&user, nil)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep)
		// ACT
		resUser, err := userServ.GetSelf(ctx)

		require.Nil(t, err)
		require.True(t, models.CmpUsers(&user, resUser))
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "GetByID", ctx, user.GetID())
	})
	t.Run("error UserID from context", func(t *testing.T) {
		ctx := context.Background()

		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("UserIDFromContext error")
		mockAuthZRep.On("UserIDFromContext", ctx).Return(uuid.Nil, expectedErr)

		mockUserRep := new(userrep.MockUserRep)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep)
		// ACT
		resUser, err := userServ.GetSelf(ctx)

		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, resUser)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
	})
	t.Run("error no userID", func(t *testing.T) {
		ctx := context.Background()

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("userRep error")
		mockUserRep.On("GetByID", ctx, user.GetID()).Return(nil, expectedErr)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep)
		// ACT
		resUser, err := userServ.GetSelf(ctx)

		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, resUser)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "GetByID", ctx, user.GetID())
	})
}

func TestUserService_ChangeSubscribeToMailing(t *testing.T) {
	userCreator := testobj.NewUserMother()
	user := userCreator.DefaultUser(uuid.New())

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		subscr := false

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("UpdateSubscribeToMailing", ctx, user.GetID(), subscr).Return(nil)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep)
		// ACT
		err := userServ.ChangeSubscribeToMailing(ctx, subscr)

		require.Nil(t, err)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "UpdateSubscribeToMailing", ctx, user.GetID(), subscr)
	})
	t.Run("error subscribe in rep", func(t *testing.T) {
		ctx := context.Background()
		subscr := false

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("userRep error")
		mockUserRep.On("UpdateSubscribeToMailing", ctx, user.GetID(), subscr).Return(expectedErr)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep)
		// ACT
		err := userServ.ChangeSubscribeToMailing(ctx, subscr)

		require.ErrorIs(t, err, expectedErr)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "UpdateSubscribeToMailing", ctx, user.GetID(), subscr)
	})
}
