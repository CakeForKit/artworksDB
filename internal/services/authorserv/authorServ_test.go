package authorserv_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/authorserv"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type AuthorServiceSuite struct {
	suite.Suite
}

func TestAuthorService(t *testing.T) {
	suite.RunSuite(t, new(AuthorServiceSuite))
}

func (s *AuthorServiceSuite) TestAuthorServ_GetAll(t provider.T) {
	authorCreator := testobj.NewAuthorMother()

	t.WithNewStep("success return 2 authors", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authors := []*models.Author{
			authorCreator.AuthorP(uuid.New()),
			authorCreator.AuthorP(uuid.New()),
		}
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("GetAll", ctx).Return(authors, nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		resAuthors, err := authorServ.GetAll(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(authors), len(resAuthors))
		for i := range resAuthors {
			sCtx.Assert().True(authors[i].Equal(resAuthors[i]))
		}
		mockAuthorRep.AssertCalled(t, "GetAll", ctx)
	}, allure.NewParameter("scenario", "success with 2 authors"))

	t.WithNewStep("success return 0 authors", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authors := []*models.Author{}
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("GetAll", ctx).Return(authors, nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		resAuthors, err := authorServ.GetAll(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Empty(resAuthors)
		mockAuthorRep.AssertCalled(t, "GetAll", ctx)
	}, allure.NewParameter("scenario", "success with empty result"))

	t.WithNewStep("error in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authors := []*models.Author{}
		expectedErr := errors.New("userRep error")
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("GetAll", ctx).Return(authors, expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		_, err := authorServ.GetAll(ctx)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockAuthorRep.AssertCalled(t, "GetAll", ctx)
	}, allure.NewParameter("scenario", "repository error"))
}

func (s *AuthorServiceSuite) TestAuthorServ_Add(t provider.T) {
	authorCreator := testobj.NewAuthorMother()

	t.WithNewStep("success add author", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		author := authorCreator.AuthorP(uuid.New())

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Add", ctx, author).Return(nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Add(ctx, author)

		sCtx.Assert().NoError(err)
		mockAuthorRep.AssertCalled(t, "Add", ctx, author)
	}, allure.NewParameter("scenario", "successful addition"))

	t.WithNewStep("error in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		author := authorCreator.AuthorP(uuid.New())
		expectedErr := errors.New("database error")

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Add", ctx, author).Return(expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Add(ctx, author)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockAuthorRep.AssertCalled(t, "Add", ctx, author)
	}, allure.NewParameter("scenario", "repository error"))
}

func (s *AuthorServiceSuite) TestAuthorServ_Update(t provider.T) {
	authorCreator := testobj.NewAuthorMother()

	t.WithNewStep("success update author", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authorID := uuid.New()
		updateReq := authorCreator.AuthorUpdateReq()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)")).Return(nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Update(ctx, authorID, updateReq)

		sCtx.Assert().NoError(err)
		mockAuthorRep.AssertCalled(t, "Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)"))
	}, allure.NewParameter("scenario", "successful update"))

	t.WithNewStep("error in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authorID := uuid.New()
		updateReq := models.AuthorUpdateReq{
			Name: "Updated Name",
		}
		expectedErr := errors.New("update error")

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)")).Return(expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Update(ctx, authorID, updateReq)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockAuthorRep.AssertCalled(t, "Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)"))
	}, allure.NewParameter("scenario", "repository error"))

	t.WithNewStep("error in update function", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authorID := uuid.New()
		updateReq := models.AuthorUpdateReq{}

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)")).
			Return(errors.New("validation error"))

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Update(ctx, authorID, updateReq)

		sCtx.Assert().Error(err)
		mockAuthorRep.AssertCalled(t, "Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)"))
	}, allure.NewParameter("scenario", "validation error"))
}

func (s *AuthorServiceSuite) TestAuthorServ_Delete(t provider.T) {
	t.WithNewStep("success delete author", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authorID := uuid.New()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(false, nil)
		mockAuthorRep.On("Delete", ctx, authorID).Return(nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)

		sCtx.Assert().NoError(err)
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertCalled(t, "Delete", ctx, authorID)
	}, allure.NewParameter("scenario", "successful deletion"))

	t.WithNewStep("error author has linked artworks", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authorID := uuid.New()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(true, nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)

		sCtx.Assert().ErrorIs(err, authorserv.ErrHasLinkedArtworks)
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertNotCalled(t, "Delete", ctx, authorID)
	}, allure.NewParameter("scenario", "has linked artworks error"))

	t.WithNewStep("error checking artworks", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authorID := uuid.New()
		expectedErr := errors.New("check artworks error")

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(false, expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)

		sCtx.Assert().Error(err)
		sCtx.Assert().Contains(err.Error(), "authorServ.Delete")
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertNotCalled(t, "Delete", ctx, authorID)
	}, allure.NewParameter("scenario", "artworks check error"))

	t.WithNewStep("error in delete", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		authorID := uuid.New()
		expectedErr := errors.New("delete error")

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(false, nil)
		mockAuthorRep.On("Delete", ctx, authorID).Return(expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertCalled(t, "Delete", ctx, authorID)
	}, allure.NewParameter("scenario", "delete operation error"))
}
