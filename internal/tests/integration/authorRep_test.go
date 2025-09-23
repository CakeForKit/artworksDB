package integration_test

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type AuthorRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx           context.Context
	authorCreator testobj.AuthorMother
	authorRep     authorrep.AuthorRep
}

func TestAuthorRepSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthorRepSuite))
}

func (s *AuthorRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.authorCreator = testobj.NewAuthorMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.authorRep, err = authorrep.NewAuthorRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *AuthorRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "author")
}

func (s *AuthorRepSuite) AfterAll(t provider.T) {
	if s.authorRep != nil {
		s.authorRep.Close()
	}
}
