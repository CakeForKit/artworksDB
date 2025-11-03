package e2eauth

import (
	"fmt"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailreader"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
)

type GodogTestConfig struct {
	BaseURL string
	AppCnfg *cnfg.AppConfig

	DBCreds         *cnfg.DatebaseCredentials
	DBCnfg          *cnfg.DatebaseConfig
	RedisCreds      *cnfg.RedisCredentials
	EmailCnfg       *cnfg.EmailCnfg
	EmailReaderCnfg *cnfg.EmailReaderCnfg
	EmailReader     *emailreader.EmailReader
}

func NewGodogTestConfig() (*GodogTestConfig, error) {
	testConf := &GodogTestConfig{}
	appCnfg, dbCreds, dbCnfg, redisCreds, err := fixtures.InitTestConfig("../../../../configs/")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	testConf.AppCnfg = appCnfg
	testConf.DBCreds = dbCreds
	testConf.DBCnfg = dbCnfg
	testConf.RedisCreds = redisCreds
	testConf.BaseURL = fmt.Sprintf("http://%s:%d/api/v1", appCnfg.AppHost, appCnfg.Port)

	testConf.EmailCnfg = cnfg.LoadEmailCnfg()
	testConf.EmailReaderCnfg = cnfg.LoadEmailReaderCnfg()
	testConf.EmailReader = emailreader.NewEmailReader(
		testConf.EmailReaderCnfg.Host,     // IMAP хост
		testConf.EmailReaderCnfg.Port,     // IMAP порт
		testConf.EmailReaderCnfg.Username, // Email
		testConf.EmailReaderCnfg.Password, // Пароль
	)

	fmt.Printf("testConf.EmailCnfg: %v\n", testConf.EmailCnfg)
	fmt.Printf("emailReaderCnfg: %v\n\n", testConf.EmailReaderCnfg)
	return testConf, nil
}
