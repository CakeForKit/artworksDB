package cnfg

type EmailCnfg struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func DefaultEmailCnfg() EmailCnfg {
	return EmailCnfg{
		Host:     "smtp.mail.ru",
		Port:     "587",
		Username: "ktestapp@mail.ru",
		Password: "t5rEDO1haK1YxUdZUrjW",
		From:     "ktestapp@mail.ru",
	}
}

type EmailReaderCnfg struct {
	Host     string
	Port     string
	Username string
	Password string
}

func DefaultEmailReaderCnfg() EmailCnfg {
	return EmailCnfg{
		Host:     "imap.mail.ru",
		Port:     "993",
		Username: "tmpforread@mail.ru",
		Password: "KouQGXYtAXiO73mBcdk6",
	}
}
