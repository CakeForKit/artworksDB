package mailing_test

// func TestMailingService_SendMailToAllUsers(t *testing.T) {
// 	userCreator := testobj.NewUserMother()
// 	user := userCreator.DefaultUser(uuid.New())

// 	t.Run("success", func(t *testing.T) {
// 		ctx := context.Background()

// 		mockUserRep := new(userrep.MockUserRep)
// 		mockUserRep.On("GetByID", ctx, user.GetID()).Return(&user, nil)

// 		mailServ := mailing.NewGmailSender(mockUserRep, "name", "fromEmailAddress", "fromEmailPassword")
// 		// ACT
// 		resUser, err := mailServ.SendMailToAllUsers(ctx)

// 		require.Nil(t, err)
// 		require.True(t, models.CmpUsers(&user, resUser))
// 		mockAuthZRep.AssertCalled(t, "UserIDFromContext", ctx)
// 		mockUserRep.AssertCalled(t, "GetByID", ctx, user.GetID())
// 	})

// }

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
