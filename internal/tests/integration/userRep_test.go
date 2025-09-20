package integration_test

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures"
	fixturesrep "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures/fixtures_rep"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type UserRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx         context.Context
	userCreator testobj.UserMother
	userRep     userrep.UserRep
}

func TestUserRepSuite(t *testing.T) {
	suite.RunSuite(t, new(UserRepSuite))
}

func (s *UserRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.userCreator = testobj.NewUserMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error
		s.userRep, err = userrep.NewUserRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err, "Cannot create UserRep")
	})
}

func (s *UserRepSuite) AfterAll(t provider.T) {
	if s.userRep != nil {
		s.userRep.Close()
	}
}

// func (s *UserRepSuite) cleanupTestData() {

// 	users, err := s.userRep.GetAll(s.ctx)
// 	if err == nil {
// 		for _, user := range users {
// 			s.userRep.Delete(s.ctx, user.GetID())
// 		}
// 	}
// }

func (s *UserRepSuite) TestUserRep_GetAll(t provider.T) {

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		users := []*models.User{
			s.userCreator.DefaultUserP(uuid.New()),
			s.userCreator.DefaultUserP(uuid.New()),
		}
		fixturesrep.AddTestUser(t, s.ctx, s.userRep, users)
		defer fixturesrep.DelTestUser(t, s.ctx, s.userRep, users)

		// ACT
		resUsers, err := s.userRep.GetAll(s.ctx)

		sCtx.Require().NoError(err)
		fixturesrep.AssertUsersAreInRes(t, users, resUsers)
	})
}
