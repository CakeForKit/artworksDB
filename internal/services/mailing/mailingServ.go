package mailing

import (
	"fmt"
	"strings"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
)

type MailingService interface {
	SendMailToAllUsers([]*models.Event) error
	SubscribeToMailing(id uuid.UUID) error
	UnSubscribeToMailing(id uuid.UUID) error
}

type mailingService struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
	userRep           userrep.UserRep
}

func NewGmailSender(urep userrep.UserRep,
	name string, fromEmailAddress string, fromEmailPassword string) MailingService {
	return &mailingService{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
		userRep:           urep,
	}
}

func (m *mailingService) SendMailToAllUsers(events []*models.Event) error {
	users := m.userRep.GetAllSubscribed()
	if len(users) == 0 {
		return nil
	}
	fmt.Printf("Сообщение отправлено пользовтелям:\n")
	msgText := m.generateMessageText(events)
	for _, u := range users {
		fmt.Printf("%s, ", u.GetMail())
	}
	fmt.Printf("\n")
	fmt.Println(msgText) // TODO to log
	return nil
}

func (m *mailingService) generateMessageText(events []*models.Event) string {
	var arre []string = make([]string, len(events)+1)
	var i int = 0
	for ; i < len(events); i++ {
		arre[i] = events[i].TextAbout()
	}
	arre[i] = fmt.Sprintf("from %s (%s)", m.name, m.fromEmailAddress)
	return strings.Join(arre, "\n")
}

func (m *mailingService) SubscribeToMailing(id uuid.UUID) error {
	updatefunc := func(u *models.User) (*models.User, error) {
		updatedUser, err := models.NewUser(
			u.GetID(),
			u.GetUsername(),
			u.GetLogin(),
			u.GetHashedPassword(),
			u.GetCreatedAt(),
			u.GetMail(),
			true,
		)
		return &updatedUser, err
	}
	_, err := m.userRep.Update(id, updatefunc)
	return err
}

func (m *mailingService) UnSubscribeToMailing(id uuid.UUID) error {
	updatefunc := func(u *models.User) (*models.User, error) {
		updatedUser, err := models.NewUser(
			u.GetID(),
			u.GetUsername(),
			u.GetLogin(),
			u.GetHashedPassword(),
			u.GetCreatedAt(),
			u.GetMail(),
			false,
		)
		return &updatedUser, err
	}
	_, err := m.userRep.Update(id, updatefunc)
	return err
}
