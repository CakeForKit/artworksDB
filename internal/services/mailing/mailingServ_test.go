package mailing_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/mailing"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type MailingServiceSuite struct {
	suite.Suite
}

func TestMailingService(t *testing.T) {
	suite.RunSuite(t, new(MailingServiceSuite))
}

func (s *MailingServiceSuite) TestMailingService_SendMailToAllUsers(t provider.T) {
	eventCreator := testobj.NewEventMother()
	userCreator := testobj.NewUserMother()

	t.WithNewStep("success with subscribed users", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{
			eventCreator.EventP(uuid.New()),
			eventCreator.EventP(uuid.New()),
		}
		users := []*models.User{
			userCreator.DefaultUserP(uuid.New()),
			userCreator.DefaultUserP(uuid.New()),
			userCreator.DefaultUserP(uuid.New()),
		}
		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAllSubscribed", ctx).Return(users, nil)

		mailServ := mailing.NewGmailSender(mockUserRep, "Gallery Name", "gallery@example.com", "password123")
		// ACT
		msgText, userIDs, err := mailServ.SendMailToAllUsers(ctx, events)

		sCtx.Require().NoError(err)
		sCtx.Assert().NotEmpty(msgText)
		sCtx.Assert().True(len(users) == len(userIDs))
		for i := range users {
			sCtx.Assert().True(users[i].GetID() == userIDs[i])
		}
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.WithNewStep("success with no subscribed users", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{
			eventCreator.EventP(uuid.New()),
		}
		users := []*models.User{}

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAllSubscribed", ctx).Return(users, nil)

		mailServ := mailing.NewGmailSender(mockUserRep, "Gallery Name", "gallery@example.com", "password123")
		// ACT
		msgText, userIDs, err := mailServ.SendMailToAllUsers(ctx, events)

		sCtx.Require().NoError(err)
		sCtx.Assert().NotEmpty(msgText)
		sCtx.Assert().True(len(userIDs) == 0)
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.WithNewStep("success with empty events", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{}
		users := []*models.User{
			userCreator.DefaultUserP(uuid.New()),
		}

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAllSubscribed", ctx).Return(users, nil)

		mailServ := mailing.NewGmailSender(mockUserRep, "Gallery Name", "gallery@example.com", "password123")
		// ACT
		msgText, userIDs, err := mailServ.SendMailToAllUsers(ctx, events)

		sCtx.Require().NoError(err)
		sCtx.Assert().NotEmpty(msgText)
		sCtx.Assert().True(len(userIDs) == 1)
		sCtx.Assert().True(users[0].GetID() == userIDs[0])
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.WithNewStep("error getting subscribed users", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{
			eventCreator.EventP(uuid.New()),
		}
		expectedErr := errors.New("database error")

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAllSubscribed", ctx).Return(nil, expectedErr)

		mailServ := mailing.NewGmailSender(mockUserRep, "Gallery Name", "gallery@example.com", "password123")
		// ACT
		msgText, userIDs, err := mailServ.SendMailToAllUsers(ctx, events)

		sCtx.Require().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		sCtx.Assert().Empty(msgText)
		sCtx.Assert().Nil(userIDs)
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.WithNewStep("success with single user", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{
			eventCreator.EventP(uuid.New()),
		}
		users := []*models.User{
			userCreator.DefaultUserP(uuid.New()),
		}

		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAllSubscribed", ctx).Return(users, nil)

		mailServ := mailing.NewGmailSender(mockUserRep, "Gallery Name", "gallery@example.com", "password123")
		// ACT
		msgText, userIDs, err := mailServ.SendMailToAllUsers(ctx, events)

		sCtx.Require().NoError(err)
		sCtx.Assert().NotEmpty(msgText)
		sCtx.Assert().True(len(userIDs) == 1)
		sCtx.Assert().True(users[0].GetID() == userIDs[0])
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})
}
