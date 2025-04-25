package mailing

import (
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
	name, email, password := createTestConfig()
	tests := []struct {
		name            string
		subscribedUsers []*models.User
		events          []*models.Event
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
			expectedError: nil,
		},
		{
			name:            "no subscribed users",
			subscribedUsers: []*models.User{},
			events:          []*models.Event{createTestEvent()},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockuserrep.MockUserRep)
			service := NewGmailSender(mockRepo, name, email, password)

			mockRepo.On("GetAllSubscribed").Return(tt.subscribedUsers)

			err := service.SendMailToAllUsers(tt.events)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMailingService_GenerateMessageText(t *testing.T) {
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
			mockRepo := new(mockuserrep.MockUserRep)
			service := NewGmailSender(mockRepo, name, email, password)

			result := service.GenerateMessageText(tt.events)
			assert.Equal(t, tt.expectedText, result)
		})
	}
}

func TestMailingService_SubscribeToMailing(t *testing.T) {
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
			mockRepo := new(mockuserrep.MockUserRep)
			service := NewGmailSender(mockRepo, name, email, password)
			mockRepo.On("UpdateSubscribeToMailing", tt.userID, true).Return(tt.mockError)

			err := service.SubscribeToMailing(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMailingService_UnSubscribeToMailing(t *testing.T) {
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
			mockRepo := new(mockuserrep.MockUserRep)
			service := NewGmailSender(mockRepo, name, email, password)
			mockRepo.On("UpdateSubscribeToMailing", tt.userID, false).Return(tt.mockError)

			err := service.UnSubscribeToMailing(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
