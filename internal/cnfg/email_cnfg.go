package cnfg

import (
	"github.com/spf13/viper"
)

type EmailCnfg struct {
	Host     string
	Port     string
	Username string
	Password string `mapstructure:"APP_EMAIL_PASSWORD"`
	From     string `mapstructure:"APP_EMAIL"`
}

func LoadEmailCnfg() *EmailCnfg {
	viper.AutomaticEnv()
	password := viper.GetString("APP_EMAIL_PASSWORD")
	from := viper.GetString("APP_EMAIL")

	return &EmailCnfg{
		Host:     "smtp.mail.ru",
		Port:     "587",
		Username: from,
		Password: password,
		From:     from,
	}
}

type EmailReaderCnfg struct {
	Host     string
	Port     string
	Username string `mapstructure:"TEST_USER_EMAIL"`
	Password string `mapstructure:"TEST_USER_EMAIL_PASSWORD"`
}

func LoadEmailReaderCnfg() *EmailReaderCnfg {
	viper.AutomaticEnv()
	password := viper.GetString("TEST_USER_EMAIL_PASSWORD")
	from := viper.GetString("TEST_USER_EMAIL")

	return &EmailReaderCnfg{
		Host:     "imap.mail.ru",
		Port:     "993",
		Username: from,
		Password: password,
	}
}
