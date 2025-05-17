package mailing

import (
	"context"
	"fmt"
	"strings"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
)

type MailingService interface {
	SendMailToAllUsers(ctx context.Context, events []*models.Event) (string, uuid.UUIDs, error)
	GenerateMessageText(ctx context.Context, events []*models.Event) string
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

func (m *mailingService) SendMailToAllUsers(ctx context.Context, events []*models.Event) (string, uuid.UUIDs, error) {
	msgText := m.GenerateMessageText(ctx, events)

	users, err := m.userRep.GetAllSubscribed(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("SendMailToAllUsers: %v", err)
	}
	if len(users) == 0 {
		return msgText, nil, nil
	}
	var userIDs uuid.UUIDs
	for _, u := range users {
		userIDs = append(userIDs, u.GetID())
	}
	return msgText, userIDs, nil
}

func (m *mailingService) GenerateMessageText(ctx context.Context, events []*models.Event) string {
	var arre []string = make([]string, len(events)+1)
	var i int = 0
	for ; i < len(events); i++ {
		arre[i] = events[i].TextAbout()
	}
	arre[i] = fmt.Sprintf("from %s (%s)", m.name, m.fromEmailAddress)
	return strings.Join(arre, "\n")
}
