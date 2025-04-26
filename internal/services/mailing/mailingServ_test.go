package mailing

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep/mockuserrep"
	"github.com/google/uuid"
	"github.com/stateio/testify/require"
)

func TestMailingService(t *testing.T) {
	validName := "museum"
	validfromEmailAddress := "museum@mail.ru"
	validfromEmailPassword := "1234"
	validEvent1, err := models.NewEvent(
		"Выставка современного искусства",
		time.Date(2023, time.November, 15, 10, 0, 0, 0, time.UTC),
		time.Date(2023, time.November, 30, 18, 0, 0, 0, time.UTC),
		"ул. Творческая, 15",
		true,
		[]*models.Artwork{},
	)
	require.NoError(t, err)
	validEvent2, err := models.NewEvent(
		"Закрытый аукцион",
		time.Date(2023, time.December, 5, 19, 0, 0, 0, time.UTC),
		time.Date(2023, time.December, 5, 22, 0, 0, 0, time.UTC),
		"ул. Коллекционная, 42",
		false,
		[]*models.Artwork{},
	)
	require.NoError(t, err)
	validUser1, err := models.NewUser(
		uuid.New(), "user1", "login1", "hash1", time.Now(), "user1@mail.ru", true,
	)
	require.NoError(t, err)
	validUser2, err := models.NewUser(
		uuid.New(), "user2", "login2", "hash2", time.Now(), "user2@mail.ru", false,
	)
	require.NoError(t, err)
	validUser3, err := models.NewUser(
		uuid.New(), "user3", "login3", "hash3", time.Now(), "user3@mail.ru", true,
	)
	require.NoError(t, err)
	validUsers := []*models.User{&validUser1, &validUser2, &validUser3}

	expectedMsgText := fmt.Sprintf("%s\n%s\nfrom %s (%s)",
		validEvent1.TextAbout(), validEvent2.TextAbout(),
		validName, validfromEmailAddress)
	validEvents := []*models.Event{
		&validEvent1,
		&validEvent2,
	}
	t.Run("SendMailToAllUsers", func(t *testing.T) {
		userRep := new(mockuserrep.MockUserRep)
		userRep.On("GetAllSubscribed").Return(validUsers)
		mailingServ := &mailingService{
			name:              validName,
			fromEmailAddress:  validfromEmailAddress,
			fromEmailPassword: validfromEmailPassword,
			userRep:           userRep,
		}
		err = mailingServ.SendMailToAllUsers(validEvents)
		require.NoError(t, err)
	})
	t.Run("generateMessageText", func(t *testing.T) {
		userRep := new(mockuserrep.MockUserRep)
		mailingServ := &mailingService{
			name:              validName,
			fromEmailAddress:  validfromEmailAddress,
			fromEmailPassword: validfromEmailPassword,
			userRep:           userRep,
		}
		msgText := mailingServ.generateMessageText(validEvents)
		require.True(t, strings.Compare(expectedMsgText, msgText) == 0)
	})
}
