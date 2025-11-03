package fixtures

import (
	"fmt"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailreader"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type BaseE2ESuite struct {
	suite.Suite
	BaseURL         string
	AppCnfg         *cnfg.AppConfig
	DBCreds         *cnfg.DatebaseCredentials
	DBCnfg          *cnfg.DatebaseConfig
	RedisCreds      *cnfg.RedisCredentials
	EmailCnfg       *cnfg.EmailCnfg
	EmailReaderCnfg *cnfg.EmailReaderCnfg
	EmailReader     *emailreader.EmailReader
}

func (s *BaseE2ESuite) BeforeAll(t provider.T) {
	t.WithNewStep("Load configuration", func(sCtx provider.StepCtx) {
		appCnfg, dbCreds, dbCnfg, redisCreds, err := InitTestConfig("../../../../configs/") // испольнение теста из: artworksDB/internal/tests/integration_test
		sCtx.Require().NoError(err, "Failed to load config")
		s.AppCnfg = appCnfg
		s.DBCreds = dbCreds
		s.DBCnfg = dbCnfg
		s.RedisCreds = redisCreds
		s.BaseURL = fmt.Sprintf("http://%s:%d/api/v1", appCnfg.AppHost, appCnfg.Port)
		s.EmailCnfg = cnfg.LoadEmailCnfg()
		s.EmailReaderCnfg = cnfg.LoadEmailReaderCnfg()
		s.EmailReader = emailreader.NewEmailReader(
			s.EmailReaderCnfg.Host,     // IMAP хост
			s.EmailReaderCnfg.Port,     // IMAP порт
			s.EmailReaderCnfg.Username, // Email
			s.EmailReaderCnfg.Password, // Пароль
		)
	})
}
