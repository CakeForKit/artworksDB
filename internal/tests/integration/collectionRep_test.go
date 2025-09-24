package integration

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures"
	fixturesrep "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures/fixtures_rep"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CollectionRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx               context.Context
	collectionCreator testobj.CollectionMother
	collectionRep     collectionrep.CollectionRep
}

func TestCollectionRepSuite(t *testing.T) {
	suite.RunSuite(t, new(CollectionRepSuite))
}

func (s *CollectionRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.collectionCreator = testobj.NewCollectionMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.collectionRep, err = collectionrep.NewCollectionRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *CollectionRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "collection")
}

func (s *CollectionRepSuite) AfterAll(t provider.T) {
	if s.collectionRep != nil {
		s.collectionRep.Close()
	}
}

func (s *CollectionRepSuite) TestCollectionRep_GetAll(t provider.T) {
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		collections := []*models.Collection{
			s.collectionCreator.CollectionP(uuid.New()),
			s.collectionCreator.CollectionP(uuid.New()),
		}
		fixturesrep.AddTestCollections(t, s.ctx, collections, s.collectionRep)
		defer fixturesrep.DelTestCollections(t, s.ctx, collections, s.collectionRep)

		// ACT
		resCollections, err := s.collectionRep.GetAll(s.ctx)

		sCtx.Require().NoError(err)
		fixturesrep.AssertCollectionsAreInRes(t, collections, resCollections)
	})
}
