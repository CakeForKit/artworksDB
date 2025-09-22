package fixtures

import (
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type BaseE2ESuite struct {
	suite.Suite
	BaseURL    string
	AppCnfg    *cnfg.AppConfig
	DBCreds    *cnfg.DatebaseCredentials
	DBCnfg     *cnfg.DatebaseConfig
	RedisCreds *cnfg.RedisCredentials
}

func (s *BaseE2ESuite) BeforeAll(t provider.T) {
	t.WithNewStep("Load configuration", func(sCtx provider.StepCtx) {
		appCnfg, dbCreds, dbCnfg, redisCreds, err := InitTestConfig("../../../../configs/") // испольнение теста из: artworksDB/internal/tests/integration_test
		sCtx.Require().NoError(err, "Failed to load config")
		s.AppCnfg = appCnfg
		s.DBCreds = dbCreds
		s.DBCnfg = dbCnfg
		s.RedisCreds = redisCreds
		s.BaseURL = fmt.Sprintf("http://localhost:%d/api/v1", appCnfg.Port)
	})
}
