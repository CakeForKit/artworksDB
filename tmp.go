package main

import (
	"fmt"
	"log"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailreader"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailserv"
)

func main() {
	emailCnfg := cnfg.LoadEmailCnfg()
	s := emailserv.NewEmailService(
		emailCnfg.Host,
		emailCnfg.Port,
		emailCnfg.Username,
		emailCnfg.Password,
		emailCnfg.From,
	)
	fmt.Printf("mail: %s, password: %s\n", emailCnfg.From, emailCnfg.Password)
	err := s.SendEmail([]string{"tmpforread@mail.ru"}, "subject", "body2")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	emailReaderCnfg := cnfg.LoadEmailReaderCnfg()

	emailReader := emailreader.NewEmailReader(
		emailReaderCnfg.Host,     // IMAP хост
		emailReaderCnfg.Port,     // IMAP порт
		emailReaderCnfg.Username, // Email
		emailReaderCnfg.Password, // Пароль
	)

	emails, err := emailReader.FindEmailByCriteria(emailreader.SearchCriteria{
		From:    emailCnfg.From,
		Subject: "subject",
	})
	if err != nil {
		log.Fatalf("Ошибка чтения писем: %v", err)
	}
	emailReader.PrintEmail(emails)
}
