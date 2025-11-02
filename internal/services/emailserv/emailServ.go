package emailserv

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewEmailService(host, port, username, password, from string) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (es *EmailService) SendEmail(to []string, subject, body string) error {
	message := []byte(
		"To: " + strings.Join(to, ",") + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n")

	auth := smtp.PlainAuth("", es.Username, es.Password, es.Host)
	address := es.Host + ":" + es.Port
	err := smtp.SendMail(address, auth, es.From, to, message)
	if err != nil {
		return fmt.Errorf("ошибка отправки email: %w", err)
	}

	return nil
}
