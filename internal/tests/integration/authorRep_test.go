package integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
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

type AuthorRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx            context.Context
	authorCreator  testobj.AuthorMother
	artworkCreator testobj.ArtworkMother
	authorRep      authorrep.AuthorRep
	artworkRep     artworkrep.ArtworkRep
	collectionRep  collectionrep.CollectionRep
}

func TestAuthorRepSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthorRepSuite))
}

func (s *AuthorRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.authorCreator = testobj.NewAuthorMother()
	s.artworkCreator = testobj.NewArtworkMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error
		s.authorRep, err = authorrep.NewAuthorRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.collectionRep, err = collectionrep.NewCollectionRep(s.ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.artworkRep, err = artworkrep.NewArtworkRep(s.ctx, s.AppCnfg.Datebase, s.DBCreds, s.DBCnfg)
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
	if s.artworkRep != nil {
		s.artworkRep.Close()
	}
	if s.collectionRep != nil {
		s.collectionRep.Close()
	}
}

func (s *AuthorRepSuite) TestAuthorRep_GetAll(t provider.T) {
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		authors := []*models.Author{
			s.authorCreator.AuthorP(uuid.New()),
			s.authorCreator.AuthorP(uuid.New()),
		}
		fixturesrep.AddTestAuthors(t, s.ctx, authors, s.authorRep)
		defer fixturesrep.DelTestAuthors(t, s.ctx, authors, s.authorRep)

		// ACT
		resAuthors, err := s.authorRep.GetAll(s.ctx)

		sCtx.Require().NoError(err)
		fixturesrep.AssertAuthorsAreInRes(t, authors, resAuthors)
	})
}

func (s *AuthorRepSuite) TestAuthorRep_GetByID(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		authors := []*models.Author{
			s.authorCreator.AuthorP(uuid.New()),
		}
		fixturesrep.AddTestAuthors(t, s.ctx, authors, s.authorRep)
		defer fixturesrep.DelTestAuthors(t, s.ctx, authors, s.authorRep)

		// ACT
		resAuthor, err := s.authorRep.GetByID(s.ctx, authors[0].GetID())

		t.Require().NoError(err)
		t.Require().True(resAuthor.Equal(authors[0]))
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.authorRep.GetByID(s.ctx, uuid.New())

		t.Require().ErrorIs(err, authorrep.ErrAuthorNotFound)
	})
}

func (s *AuthorRepSuite) TestAuthorRep_Add(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		author := s.authorCreator.AuthorP(uuid.New())
		defer fixturesrep.DelTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)

		// ACT
		err := s.authorRep.Add(s.ctx, author)

		t.Require().NoError(err)

		// Verify author was added
		resAuthor, err := s.authorRep.GetByID(s.ctx, author.GetID())
		t.Require().NoError(err)
		t.Require().True(resAuthor.Equal(author))
	})

	t.Run("Duplicate error", func(t provider.T) {
		author := s.authorCreator.AuthorP(uuid.New())
		fixturesrep.AddTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)
		defer fixturesrep.DelTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)

		// ACT - try to add author with same ID
		err := s.authorRep.Add(s.ctx, author)

		t.Require().Error(err)
	})
}

func (s *AuthorRepSuite) TestAuthorRep_Delete(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		author := s.authorCreator.AuthorP(uuid.New())
		fixturesrep.AddTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)

		// ACT
		err := s.authorRep.Delete(s.ctx, author.GetID())

		t.Require().NoError(err)

		// Verify author was deleted
		_, err = s.authorRep.GetByID(s.ctx, author.GetID())
		t.Require().ErrorIs(err, authorrep.ErrAuthorNotFound)
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.authorRep.Delete(s.ctx, uuid.New())

		t.Require().Error(err)
	})
}

func (s *AuthorRepSuite) TestAuthorRep_Update(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		author := s.authorCreator.AuthorP(uuid.New())
		fixturesrep.AddTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)
		defer fixturesrep.DelTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)

		newAuthor, err := models.NewAuthor(
			author.GetID(),
			"Updated Author Name",
			author.GetBirthYear()-1,
			author.GetDeathYear()+10,
		)
		t.Require().NoError(err)
		funcUpdate := func(a *models.Author) (*models.Author, error) {
			return &newAuthor, nil
		}

		// ACT
		err = s.authorRep.Update(s.ctx, author.GetID(), funcUpdate)

		t.Require().NoError(err)

		// Verify changes persisted
		updatedAuthor, err := s.authorRep.GetByID(s.ctx, author.GetID())
		t.Require().NoError(err)
		t.Require().True(newAuthor.Equal(updatedAuthor))
	})

	t.Run("Not found", func(t provider.T) {
		funcUpdate := func(a *models.Author) (*models.Author, error) {
			return a, nil
		}

		// ACT
		err := s.authorRep.Update(s.ctx, uuid.New(), funcUpdate)

		t.Require().ErrorIs(err, authorrep.ErrAuthorNotFound)
	})

	t.Run("Update function returns error", func(t provider.T) {
		author := s.authorCreator.AuthorP(uuid.New())
		fixturesrep.AddTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)
		defer fixturesrep.DelTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)

		expectedErr := errors.New("update function error")
		funcUpdate := func(a *models.Author) (*models.Author, error) {
			return nil, expectedErr
		}

		// ACT
		err := s.authorRep.Update(s.ctx, author.GetID(), funcUpdate)

		t.Require().ErrorIs(err, expectedErr)
	})
}

func (s *AuthorRepSuite) TestAuthorRep_HasArtworks(t provider.T) {
	// t.Parallel()

	t.Run("Has artworks - true", func(t provider.T) {
		author := s.authorCreator.AuthorP(uuid.New())
		artworks := []*models.Artwork{
			s.artworkCreator.ArtworkWithAuthorP(uuid.New(), author),
			s.artworkCreator.ArtworkWithAuthorP(uuid.New(), author),
		}
		fixturesrep.AddTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, artworks, s.artworkRep, s.authorRep, s.collectionRep)
		// ACT

		hasArtworks, err := s.authorRep.HasArtworks(s.ctx, author.GetID())

		t.Require().NoError(err)
		t.Require().True(hasArtworks)
	})

	t.Run("Has artworks - false", func(t provider.T) {
		author := s.authorCreator.AuthorP(uuid.New())
		fixturesrep.AddTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)
		defer fixturesrep.DelTestAuthors(t, s.ctx, []*models.Author{author}, s.authorRep)

		// ACT
		hasArtworks, err := s.authorRep.HasArtworks(s.ctx, author.GetID())

		t.Require().NoError(err)
		t.Require().False(hasArtworks)
	})

	t.Run("Author not found", func(t provider.T) {
		// ACT
		_, err := s.authorRep.HasArtworks(s.ctx, uuid.New())

		t.Require().ErrorIs(err, authorrep.ErrAuthorNotFound)
	})
}
