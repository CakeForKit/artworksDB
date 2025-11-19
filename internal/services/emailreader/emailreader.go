package emailreader

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type EmailReader struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewEmailReader(host, port, username, password string) *EmailReader {
	return &EmailReader{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

type Email struct {
	ID      uint32
	From    string
	Subject string
	Date    string
	Body    string
}

type SearchCriteria struct {
	From    string // Email отправителя
	Subject string // Тема письма (может быть частью темы)
}

// FindEmailByCriteria ищет письмо по критериям и возвращает САМОЕ ПОСЛЕДНЕЕ найденное
func (er *EmailReader) FindEmailByCriteria(criteria SearchCriteria) (*Email, error) {
	c, err := er.connectAndLogin()
	if err != nil {
		return nil, err
	}
	defer er.logout(c)

	mailbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("ошибка выбора папки: %w", err)
	}
	if mailbox.Messages == 0 {
		return nil, fmt.Errorf("в папке INBOX нет писем")
	}

	foundEmails, err := er.fetchAndFilterEmails(c, mailbox.Messages, criteria)
	if err != nil {
		return nil, err
	}

	return er.findLatestEmail(foundEmails)
}

func (er *EmailReader) connectAndLogin() (*client.Client, error) {
	imapAddr := er.Host + ":" + er.Port
	c, err := client.DialTLS(imapAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}

	if err := c.Login(er.Username, er.Password); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", err)
	}

	return c, nil
}

func (er *EmailReader) logout(c *client.Client) {
	if err := c.Logout(); err != nil {
		log.Printf("Logout error: %v", err)
	}
}

func (er *EmailReader) fetchAndFilterEmails(c *client.Client, totalMessages uint32,
	criteria SearchCriteria) ([]*Email, error) {
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, totalMessages)

	section := &imap.BodySectionName{}
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{
			imap.FetchEnvelope,
			imap.FetchFlags,
			imap.FetchInternalDate,
			section.FetchItem(),
		}, messages)
	}()

	var foundEmails []*Email
	for msg := range messages {
		email, ok := er.processMessage(msg, section, criteria)
		if ok {
			foundEmails = append(foundEmails, email)
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("ошибка получения писем: %w", err)
	}

	return foundEmails, nil
}

func (er *EmailReader) processMessage(msg *imap.Message, section *imap.BodySectionName,
	criteria SearchCriteria) (*Email, bool) {
	fromAddress := msg.Envelope.From[0].Address()
	subject := msg.Envelope.Subject

	if !er.matchesCriteria(fromAddress, subject, criteria) {
		return nil, false
	}

	email := &Email{
		ID:      msg.SeqNum,
		From:    fromAddress,
		Subject: subject,
		Date:    msg.Envelope.Date.Format("2006-01-02 15:04"),
	}

	if msg.GetBody(section) != nil {
		if body, err := er.extractBody(msg.GetBody(section)); err == nil {
			email.Body = body
		}
	}

	return email, true
}

func (er *EmailReader) findLatestEmail(emails []*Email) (*Email, error) {
	if len(emails) == 0 {
		return nil, fmt.Errorf("письмо не найдено")
	}

	var latestEmail *Email
	for _, email := range emails {
		if latestEmail == nil || email.ID > latestEmail.ID {
			latestEmail = email
		}
	}

	return latestEmail, nil
}

func (er *EmailReader) matchesCriteria(from, subject string, criteria SearchCriteria) bool {
	if criteria.From != "" && from != criteria.From {
		return false
	}
	if criteria.Subject != "" {
		if !containsIgnoreCase(subject, criteria.Subject) {
			return false
		}
	}
	return true
}

func (er *EmailReader) extractBody(body io.Reader) (string, error) {
	mr, err := mail.CreateReader(body)
	if err != nil {
		return "", err
	}
	defer mr.Close()

	var textBody string

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			if contentType == "text/plain" {
				bodyBytes, err := io.ReadAll(p.Body)
				if err == nil {
					textBody = string(bodyBytes)
				}
			}
		}
	}

	return textBody, nil
}

func (er *EmailReader) PrintEmail(email *Email) {
	if email == nil {
		fmt.Println("❌ Письмо не найдено")
		return
	}
	fmt.Println("=====================================")
	fmt.Printf("ID: %d\n", email.ID)
	fmt.Printf("От: %s\n", email.From)
	fmt.Printf("Тема: %s\n", email.Subject)
	fmt.Printf("Дата: %s\n", email.Date)
	fmt.Printf("Тело: %s\n", email.Body)
	fmt.Println("=====================================")
}

func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}

	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}
