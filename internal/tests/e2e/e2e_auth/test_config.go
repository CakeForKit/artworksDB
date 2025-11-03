package e2eauth

import (
	"fmt"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailreader"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
)

type GodogTestConfig struct {
	BaseURL     string
	AppCnfg     *cnfg.AppConfig
	EmailCnfg   *cnfg.EmailCnfg
	DBCreds     *cnfg.DatebaseCredentials
	DBCnfg      *cnfg.DatebaseConfig
	RedisCreds  *cnfg.RedisCredentials
	EmailReader *emailreader.EmailReader
}

func NewGodogTestConfig() (*GodogTestConfig, error) {
	testConf := &GodogTestConfig{}
	appCnfg, dbCreds, dbCnfg, redisCreds, err := fixtures.InitTestConfig("../../../../configs/")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	testConf.AppCnfg = appCnfg
	testConf.EmailCnfg = cnfg.LoadEmailCnfg()

	testConf.DBCreds = dbCreds
	testConf.DBCnfg = dbCnfg
	testConf.RedisCreds = redisCreds
	testConf.BaseURL = fmt.Sprintf("http://%s:%d/api/v1", appCnfg.AppHost, appCnfg.Port)
	emailReaderCnfg := cnfg.LoadEmailReaderCnfg()
	testConf.EmailReader = emailreader.NewEmailReader(
		emailReaderCnfg.Host,     // IMAP хост
		emailReaderCnfg.Port,     // IMAP порт
		emailReaderCnfg.Username, // Email
		emailReaderCnfg.Password, // Пароль
	)
	fmt.Printf("testConf.EmailCnfg: %v\n", testConf.EmailCnfg)
	fmt.Printf("emailReaderCnfg: %v\n\n", emailReaderCnfg)
	return testConf, nil
}
