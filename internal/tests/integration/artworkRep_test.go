package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/repository/artworkrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
	fixturesrep "github.com/CakeForKit/artworksDB.git/internal/tests/fixtures/fixtures_rep"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ArtworkRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx            context.Context
	artworkCreator testobj.ArtworkMother
	artworkRep     artworkrep.ArtworkRep
	authorRep      authorrep.AuthorRep
	collectionRep  collectionrep.CollectionRep
}

func TestArtworkRepSuite(t *testing.T) {
	suite.RunSuite(t, new(ArtworkRepSuite))
}

func (s *ArtworkRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.artworkCreator = testobj.NewArtworkMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.artworkRep, err = artworkrep.NewArtworkRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.authorRep, err = authorrep.NewAuthorRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.collectionRep, err = collectionrep.NewCollectionRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *ArtworkRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "artwork")
}

func (s *ArtworkRepSuite) AfterAll(t provider.T) {
	if s.artworkRep != nil {
		s.artworkRep.Close()
	}
}

func (s *ArtworkRepSuite) TestArtworkRep_GetAllArtworks(t provider.T) {
	t.Parallel()

	t.Run("Success with empty filter", func(t provider.T) {
		artworks := []*models.Artwork{
			s.artworkCreator.ArtworkP(uuid.New()),
			s.artworkCreator.ArtworkP(uuid.New()),
		}
		fixturesrep.AddTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)

		filterOps := s.artworkCreator.ArtworkFilterEmpty()
		sortOps := s.artworkCreator.ArtworkSortOps()

		// ACT
		resArtworks, err := s.artworkRep.GetAllArtworks(s.ctx, filterOps, sortOps)

		t.Require().NoError(err)
		fixturesrep.AssertArtworksAreInRes(t, artworks, resArtworks)
	})

	t.Run("Success with title filter", func(t provider.T) {
		artwork1 := s.artworkCreator.ArtworkP(uuid.New())
		artwork2 := s.artworkCreator.ArtworkP(uuid.New())
		artworks := []*models.Artwork{artwork1, artwork2}
		fixturesrep.AddTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)

		filterOps := &jsonreqresp.ArtworkFilter{
			Title:      artwork1.GetTitle(),
			AuthorName: "",
			Collection: "",
			EventID:    uuid.Nil,
		}
		sortOps := s.artworkCreator.ArtworkSortOps()

		// ACT
		resArtworks, err := s.artworkRep.GetAllArtworks(s.ctx, filterOps, sortOps)

		t.Require().NoError(err)
		fixturesrep.AssertArtworksAreInRes(t, []*models.Artwork{artwork1}, resArtworks)
		for _, v := range resArtworks {
			t.Assert().True(v.GetTitle() == filterOps.Title)
		}
	})
}

func (s *ArtworkRepSuite) TestArtworkRep_GetByID(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		artworks := []*models.Artwork{
			s.artworkCreator.ArtworkP(uuid.New()),
		}
		fixturesrep.AddTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		resArtwork, err := s.artworkRep.GetByID(s.ctx, artworks[0].GetID())

		t.Require().NoError(err)
		t.Require().True(resArtwork.Equal(artworks[0]))
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.artworkRep.GetByID(s.ctx, uuid.New())

		t.Require().ErrorIs(err, artworkrep.ErrArtworkNotFound)
	})
}

func (s *ArtworkRepSuite) TestArtworkRep_Add(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		fixturesrep.AddTestAuthors(t, s.ctx, []*models.Author{artwork.GetAuthor()}, s.authorRep)
		fixturesrep.AddTestCollections(t, s.ctx, []*models.Collection{artwork.GetCollection()}, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err := s.artworkRep.Add(s.ctx, artwork)

		t.Require().NoError(err)

		resArtwork, err := s.artworkRep.GetByID(s.ctx, artwork.GetID())
		t.Require().NoError(err)
		t.Require().True(resArtwork.Equal(artwork))
	})

	t.Run("Duplicate error", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT - try to add artwork with same ID
		err := s.artworkRep.Add(s.ctx, artwork)

		t.Require().Error(err)
	})

	t.Run("Author not found error", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		// Don't add author to repository

		// ACT
		err := s.artworkRep.Add(s.ctx, artwork)

		t.Require().Error(err)
	})

	t.Run("Collection not found error", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		// Add author but not collection
		fixturesrep.AddTestAuthors(t, s.ctx, []*models.Author{artwork.GetAuthor()}, s.authorRep)
		defer fixturesrep.DelTestAuthors(t, s.ctx, []*models.Author{artwork.GetAuthor()}, s.authorRep)

		// ACT
		err := s.artworkRep.Add(s.ctx, artwork)

		t.Require().Error(err)
	})
}

func (s *ArtworkRepSuite) TestArtworkRep_Delete(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err := s.artworkRep.Delete(s.ctx, artwork.GetID())

		t.Require().NoError(err)

		// Verify artwork was deleted
		_, err = s.artworkRep.GetByID(s.ctx, artwork.GetID())
		t.Require().ErrorIs(err, artworkrep.ErrArtworkNotFound)
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.artworkRep.Delete(s.ctx, uuid.New())

		t.Require().Error(err)
	})
}

func (s *ArtworkRepSuite) TestArtworkRep_Update(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		newArtwork, err := models.NewArtwork(
			artwork.GetID(),
			"Updated Artwork Title",
			"Updated technic",
			"Updated material",
			"Updated size",
			artwork.GetCreationYear(),
			artwork.GetAuthor(),
			artwork.GetCollection(),
		)
		t.Require().NoError(err)
		funcUpdate := func(a *models.Artwork) (*models.Artwork, error) {
			return &newArtwork, nil
		}

		// ACT
		err = s.artworkRep.Update(s.ctx, artwork.GetID(), funcUpdate)

		t.Require().NoError(err)

		updatedArtwork, err := s.artworkRep.GetByID(s.ctx, artwork.GetID())
		t.Require().NoError(err)
		t.Require().True(updatedArtwork.Equal(&newArtwork))
	})

	t.Run("Not found", func(t provider.T) {
		funcUpdate := func(a *models.Artwork) (*models.Artwork, error) {
			return a, nil
		}

		// ACT
		err := s.artworkRep.Update(s.ctx, uuid.New(), funcUpdate)

		t.Require().ErrorIs(err, artworkrep.ErrArtworkNotFound)
	})

	t.Run("Update function returns error", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		expectedErr := errors.New("update function error")
		funcUpdate := func(a *models.Artwork) (*models.Artwork, error) {
			return nil, expectedErr
		}

		// ACT
		err := s.artworkRep.Update(s.ctx, artwork.GetID(), funcUpdate)

		t.Require().ErrorIs(err, expectedErr)
	})

}
