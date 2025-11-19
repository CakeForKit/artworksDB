package fixtures

import (
	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type BaseIntegrationSuite struct {
	suite.Suite
	AppCnfg    *cnfg.AppConfig
	DBCreds    *cnfg.DatebaseCredentials
	DBCnfg     *cnfg.DatebaseConfig
	RedisCreds *cnfg.RedisCredentials
}

func (s *BaseIntegrationSuite) BeforeAll(t provider.T) {
	t.Tags("integration")
	t.WithNewStep("Load configuration", func(sCtx provider.StepCtx) {
		// испольнение теста из: artworksDB/internal/tests/integration_test
		appCnfg, dbCreds, dbCnfg, redisCreds, err := InitTestConfig("../../../configs/")
		sCtx.Require().NoError(err, "Failed to load config")
		s.AppCnfg = appCnfg
		s.DBCreds = dbCreds
		s.DBCnfg = dbCnfg
		s.RedisCreds = redisCreds
	})
}
