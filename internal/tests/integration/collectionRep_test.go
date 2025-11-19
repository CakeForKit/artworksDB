package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
	fixturesrep "github.com/CakeForKit/artworksDB.git/internal/tests/fixtures/fixtures_rep"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
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
		var err error
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

func (s *CollectionRepSuite) TestCollectionRep_GetByID(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		collections := []*models.Collection{
			s.collectionCreator.CollectionP(uuid.New()),
		}
		fixturesrep.AddTestCollections(t, s.ctx, collections, s.collectionRep)
		defer fixturesrep.DelTestCollections(t, s.ctx, collections, s.collectionRep)

		// ACT
		resCollection, err := s.collectionRep.GetByID(s.ctx, collections[0].GetID())

		t.Require().NoError(err)
		t.Require().True(resCollection.Equal(collections[0]))
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.collectionRep.GetByID(s.ctx, uuid.New())

		t.Require().ErrorIs(err, collectionrep.ErrCollectionNotFound)
	})
}

func (s *CollectionRepSuite) TestCollectionRep_Add(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		collection := s.collectionCreator.CollectionP(uuid.New())
		defer fixturesrep.DelTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)

		// ACT
		err := s.collectionRep.Add(s.ctx, collection)

		t.Require().NoError(err)

		// Verify collection was added
		resCollection, err := s.collectionRep.GetByID(s.ctx, collection.GetID())
		t.Require().NoError(err)
		t.Require().True(resCollection.Equal(collection))
	})

	t.Run("Duplicate error", func(t provider.T) {
		collection := s.collectionCreator.CollectionP(uuid.New())
		fixturesrep.AddTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)
		defer fixturesrep.DelTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)

		// ACT - try to add collection with same ID
		err := s.collectionRep.Add(s.ctx, collection)

		t.Require().Error(err)
	})
}

func (s *CollectionRepSuite) TestCollectionRep_Delete(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		collection := s.collectionCreator.CollectionP(uuid.New())
		fixturesrep.AddTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)

		// ACT
		err := s.collectionRep.Delete(s.ctx, collection.GetID())

		t.Require().NoError(err)

		// Verify collection was deleted
		_, err = s.collectionRep.GetByID(s.ctx, collection.GetID())
		t.Require().ErrorIs(err, collectionrep.ErrCollectionNotFound)
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.collectionRep.Delete(s.ctx, uuid.New())

		t.Require().Error(err)
	})
}

func (s *CollectionRepSuite) TestCollectionRep_Update(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		collection := s.collectionCreator.CollectionP(uuid.New())
		fixturesrep.AddTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)
		defer fixturesrep.DelTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)

		newCollections, err := models.NewCollection(
			collection.GetID(),
			"Updated Collection Title",
		)
		t.Require().NoError(err)
		funcUpdate := func(c *models.Collection) (*models.Collection, error) {
			return &newCollections, nil
		}

		// ACT
		err = s.collectionRep.Update(s.ctx, collection.GetID(), funcUpdate)

		t.Require().NoError(err)

		// Verify changes persisted
		updatedCollection, err := s.collectionRep.GetByID(s.ctx, collection.GetID())
		t.Require().NoError(err)
		t.Require().True(newCollections.Equal(updatedCollection))
	})

	t.Run("Not found", func(t provider.T) {
		funcUpdate := func(c *models.Collection) (*models.Collection, error) {
			return c, nil
		}

		// ACT
		err := s.collectionRep.Update(s.ctx, uuid.New(), funcUpdate)

		t.Require().ErrorIs(err, collectionrep.ErrCollectionNotFound)
	})

	t.Run("Update function returns error", func(t provider.T) {
		collection := s.collectionCreator.CollectionP(uuid.New())
		fixturesrep.AddTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)
		defer fixturesrep.DelTestCollections(t, s.ctx, []*models.Collection{collection}, s.collectionRep)

		expectedErr := errors.New("update function error")
		funcUpdate := func(c *models.Collection) (*models.Collection, error) {
			return nil, expectedErr
		}

		// ACT
		err := s.collectionRep.Update(s.ctx, collection.GetID(), funcUpdate)

		t.Require().ErrorIs(err, expectedErr)
	})
}
