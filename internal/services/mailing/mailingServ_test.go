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
	"github.com/stretchr/testify/require"
)

func TestMailingService_SendMailToAllUsers(t *testing.T) {
	eventCreator := testobj.NewEventMother()
	userCreator := testobj.NewUserMother()

	t.Run("success with subscribed users", func(t *testing.T) {
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

		require.Nil(t, err)
		require.NotEmpty(t, msgText)
		require.True(t, len(users) == len(userIDs))
		for i := range users {
			require.True(t, users[i].GetID() == userIDs[i])
		}
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.Run("success with no subscribed users", func(t *testing.T) {
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

		require.Nil(t, err)
		require.NotEmpty(t, msgText)
		require.True(t, len(userIDs) == 0)
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.Run("success with empty events", func(t *testing.T) {
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

		require.Nil(t, err)
		require.NotEmpty(t, msgText)
		require.True(t, len(userIDs) == 1)
		require.True(t, users[0].GetID() == userIDs[0])
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.Run("error getting subscribed users", func(t *testing.T) {
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

		require.Error(t, err)
		require.Contains(t, err.Error(), "SendMailToAllUsers")
		require.Empty(t, msgText)
		require.Nil(t, userIDs)
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})

	t.Run("success with single user", func(t *testing.T) {
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

		require.Nil(t, err)
		require.NotEmpty(t, msgText)
		require.True(t, len(userIDs) == 1)
		require.True(t, users[0].GetID() == userIDs[0])
		mockUserRep.AssertCalled(t, "GetAllSubscribed", ctx)
	})
}

/*
func createTestConfig() (string, string, string) {
	return "Test Museum", "museum@test.com", "test-password"
}

func createTestUser(subscribed bool) *models.User {
	user, _ := models.NewUser(
		uuid.New(),
		"test-user",
		"test-login",
		"hashed-password",
		time.Now(),
		"user@test.com",
		subscribed,
	)
	return &user
}

func createTestEvent() *models.Event {
	event, _ := models.NewEvent(
		uuid.New(),
		"Test Event",
		time.Now(),
		time.Now().Add(24*time.Hour),
		"Test Address",
		true,
		uuid.New(),
		100,
		true,
		make(uuid.UUIDs, 0),
	)
	return &event
}

func TestMailingService_SendMailToAllUsers(t *testing.T) {
	ctx := context.Background()
	name, email, password := createTestConfig()

	tests := []struct {
		name            string
		subscribedUsers []*models.User
		events          []*models.Event
		mockError       error
		expectedError   error
		expectedIDsLen  int
	}{
		{
			name: "with subscribed users",
			subscribedUsers: []*models.User{
				createTestUser(true),
				createTestUser(true),
			},
			events: []*models.Event{
				createTestEvent(),
				createTestEvent(),
			},
			mockError:      nil,
			expectedError:  nil,
			expectedIDsLen: 2,
		},
		{
			name:            "no subscribed users",
			subscribedUsers: []*models.User{},
			events:          []*models.Event{createTestEvent()},
			mockError:       nil,
			expectedError:   nil,
			expectedIDsLen:  0,
		},
		{
			name:            "repository error",
			subscribedUsers: nil,
			events:          []*models.Event{createTestEvent()},
			mockError:       assert.AnError,
			expectedError:   assert.AnError,
			expectedIDsLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRep := new(userrep.MockUserRep)
			service := NewGmailSender(userRep, name, email, password)

			userRep.On("GetAllSubscribed", ctx).Return(tt.subscribedUsers, tt.mockError)

			msgText, userIDs, err := service.SendMailToAllUsers(ctx, tt.events)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, msgText)
				assert.Equal(t, tt.expectedIDsLen, len(userIDs))
			}
			userRep.AssertExpectations(t)
		})
	}
}
*/
