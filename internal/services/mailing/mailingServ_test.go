package mailing

import (
	"context"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep/mockuserrep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
		10,
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
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:            "no subscribed users",
			subscribedUsers: []*models.User{},
			events:          []*models.Event{createTestEvent()},
			mockError:       nil,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRep := new(mockuserrep.MockUserRep)
			service := NewGmailSender(userRep, name, email, password)

			userRep.On("GetAllSubscribed", ctx).Return(tt.subscribedUsers, tt.mockError)

			err := service.SendMailToAllUsers(ctx, tt.events)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			userRep.AssertExpectations(t)
		})
	}
}

func TestMailingService_GenerateMessageText(t *testing.T) {
	ctx := context.Background()
	eventCur1 := createTestEvent()
	eventCur2 := createTestEvent()
	name, email, password := createTestConfig()
	tests := []struct {
		name         string
		events       []*models.Event
		expectedText string
	}{
		{
			name:         "single event",
			events:       []*models.Event{eventCur1},
			expectedText: eventCur1.TextAbout() + "\nfrom " + name + " (" + email + ")",
		},
		{
			name: "multiple events",
			events: []*models.Event{
				eventCur1,
				eventCur2,
			},
			expectedText: eventCur1.TextAbout() + "\n" + eventCur2.TextAbout() + "\nfrom " + name + " (" + email + ")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRep := new(mockuserrep.MockUserRep)
			service := NewGmailSender(userRep, name, email, password)

			result := service.GenerateMessageText(ctx, tt.events)
			assert.Equal(t, tt.expectedText, result)
		})
	}
}

func TestMailingService_SubscribeToMailing(t *testing.T) {
	ctx := context.Background()
	name, email, password := createTestConfig()
	tests := []struct {
		name          string
		userID        uuid.UUID
		mockError     error
		expectedError error
	}{
		{
			name:          "success",
			userID:        uuid.New(),
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "user not found",
			userID:        uuid.New(),
			mockError:     userrep.ErrUserNotFound,
			expectedError: userrep.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRep := new(mockuserrep.MockUserRep)
			service := NewGmailSender(userRep, name, email, password)
			userRep.On("UpdateSubscribeToMailing", ctx, tt.userID, true).Return(tt.mockError)

			err := service.SubscribeToMailing(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			userRep.AssertExpectations(t)
		})
	}
}

func TestMailingService_UnSubscribeToMailing(t *testing.T) {
	ctx := context.Background()
	name, email, password := createTestConfig()
	tests := []struct {
		name          string
		userID        uuid.UUID
		mockError     error
		expectedError error
	}{
		{
			name:          "success",
			userID:        uuid.New(),
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "user not found",
			userID:        uuid.New(),
			mockError:     userrep.ErrUserNotFound,
			expectedError: userrep.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRep := new(mockuserrep.MockUserRep)
			service := NewGmailSender(userRep, name, email, password)
			userRep.On("UpdateSubscribeToMailing", ctx, tt.userID, false).Return(tt.mockError)

			err := service.UnSubscribeToMailing(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			userRep.AssertExpectations(t)
		})
	}
}
