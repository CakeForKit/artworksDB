package integration

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures"
	fixturesrep "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures/fixtures_rep"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type EventRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx           context.Context
	eventCreator  testobj.EventMother
	eventRep      eventrep.EventRep
	employeeRep   employeerep.EmployeeRep
	adminRep      adminrep.AdminRep
	artworkRep    artworkrep.ArtworkRep
	authorRep     authorrep.AuthorRep
	collectionRep collectionrep.CollectionRep
}

func TestEventRepSuite(t *testing.T) {
	suite.RunSuite(t, new(EventRepSuite))
}

func (s *EventRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.eventCreator = testobj.NewEventMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.eventRep, err = eventrep.NewEventRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.artworkRep, err = artworkrep.NewArtworkRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.authorRep, err = authorrep.NewAuthorRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.collectionRep, err = collectionrep.NewCollectionRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.employeeRep, err = employeerep.NewEmployeeRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.adminRep, err = adminrep.NewAdminRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *EventRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "event")
}

func (s *EventRepSuite) AfterAll(t provider.T) {
	if s.eventRep != nil {
		s.eventRep.Close()
	}
}

func (s *EventRepSuite) TestEventRep_GetAllEvents(t provider.T) {
	t.Parallel()

	t.Run("Success with empty filter", func(t provider.T) {
		events := []*models.Event{
			s.eventCreator.EventP(uuid.New()),
			s.eventCreator.EventP(uuid.New()),
		}
		fixturesrep.AddTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		filterOps := s.eventCreator.EventFilterEmpty()

		// ACT
		resEvents, err := s.eventRep.GetAll(s.ctx, filterOps)

		t.Require().NoError(err)
		fixturesrep.AssertEventsAreInRes(t, events, resEvents)
	})
}
