package userservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/userservice"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type UserServiceSuite struct {
	suite.Suite
}

func TestUserService(t *testing.T) {
	suite.RunSuite(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) TestUserService_GetSelf(t provider.T) {
	userCreator := testobj.NewUserMother()
	user := userCreator.DefaultUser(uuid.New())
	hasher, err := hasher.NewHasher()
	t.Require().NoError(err)

	t.WithNewStep("WithNewStep", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetByID", ctx, user.GetID()).Return(&user, nil)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep, hasher)
		// ACT
		resUser, err := userServ.GetSelf(ctx)

		sCtx.Require().NoError(err)
		// sCtx.Require().Error(err)
		sCtx.Assert().True(models.CmpUsers(&user, resUser))
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "GetByID", ctx, user.GetID())
	})
	t.WithNewStep("error UserID from context", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("UserIDFromContext error")
		mockAuthZRep.On("UserIDFromContext", ctx).Return(uuid.Nil, expectedErr)

		mockUserRep := new(userrep.MockUserRep)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep, hasher)
		// ACT
		resUser, err := userServ.GetSelf(ctx)

		sCtx.Assert().ErrorIs(err, expectedErr)
		sCtx.Assert().True(resUser == nil)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
	})
	t.WithNewStep("error no userID", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("userRep error")
		mockUserRep.On("GetByID", ctx, user.GetID()).Return(nil, expectedErr)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep, hasher)
		// ACT
		resUser, err := userServ.GetSelf(ctx)

		sCtx.Assert().ErrorIs(err, expectedErr)
		sCtx.Assert().True(resUser == nil)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "GetByID", ctx, user.GetID())
	})
}

func (s *UserServiceSuite) TestUserService_ChangeSubscribeToMailing(t provider.T) {
	userCreator := testobj.NewUserMother()
	user := userCreator.DefaultUser(uuid.New())
	hasher, err := hasher.NewHasher()
	t.Require().NoError(err)

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		subscr := false

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("UpdateSubscribeToMailing", ctx, user.GetID(), subscr).Return(nil)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep, hasher)
		// ACT
		err := userServ.ChangeSubscribeToMailing(ctx, subscr)

		sCtx.Require().NoError(err)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "UpdateSubscribeToMailing", ctx, user.GetID(), subscr)
	})
	t.WithNewStep("error subscribe in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		subscr := false

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("UserIDFromContext", ctx).Return(user.GetID(), nil)

		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("userRep error")
		mockUserRep.On("UpdateSubscribeToMailing", ctx, user.GetID(), subscr).Return(expectedErr)

		userServ := userservice.NewUserService(mockUserRep, mockAuthZRep, hasher)
		// ACT
		err := userServ.ChangeSubscribeToMailing(ctx, subscr)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
		mockUserRep.AssertCalled(t, "UpdateSubscribeToMailing", ctx, user.GetID(), subscr)
	})
}
