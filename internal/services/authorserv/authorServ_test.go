package authorserv_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/authorserv"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthorServ_GetAll(t *testing.T) {
	authorCreator := testobj.NewAuthorMother()

	t.Run("success return 2 authors", func(t *testing.T) {
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
		require.Nil(t, err)
		require.True(t, len(authors) == len(resAuthors))
		for i := range len(resAuthors) {
			require.True(t, authors[i].Equals(resAuthors[i]))
		}
		mockAuthorRep.AssertCalled(t, "GetAll", ctx)
	})
	t.Run("success return 0 authors", func(t *testing.T) {
		ctx := context.Background()
		authors := []*models.Author{}
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("GetAll", ctx).Return(authors, nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		resAuthors, err := authorServ.GetAll(ctx)
		require.Nil(t, err)
		require.True(t, len(resAuthors) == 0)
		mockAuthorRep.AssertCalled(t, "GetAll", ctx)
	})
	t.Run("error in rep", func(t *testing.T) {
		ctx := context.Background()
		authors := []*models.Author{}
		expectedErr := errors.New("userRep error")
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("GetAll", ctx).Return(authors, expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		_, err := authorServ.GetAll(ctx)
		require.ErrorIs(t, err, expectedErr)
		mockAuthorRep.AssertCalled(t, "GetAll", ctx)
	})
}

func TestAuthorServ_Add(t *testing.T) {
	authorCreator := testobj.NewAuthorMother()

	t.Run("success add author", func(t *testing.T) {
		ctx := context.Background()
		author := authorCreator.AuthorP(uuid.New())

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Add", ctx, author).Return(nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Add(ctx, author)
		require.Nil(t, err)
		mockAuthorRep.AssertCalled(t, "Add", ctx, author)
	})

	t.Run("error in rep", func(t *testing.T) {
		ctx := context.Background()
		author := authorCreator.AuthorP(uuid.New())
		expectedErr := errors.New("database error")

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Add", ctx, author).Return(expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Add(ctx, author)
		require.ErrorIs(t, err, expectedErr)
		mockAuthorRep.AssertCalled(t, "Add", ctx, author)
	})
}

func TestAuthorServ_Update(t *testing.T) {
	authorCreator := testobj.NewAuthorMother()

	t.Run("success update author", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()
		updateReq := authorCreator.AuthorUpdateReq()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)")).Return(nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Update(ctx, authorID, updateReq)
		require.Nil(t, err)
		mockAuthorRep.AssertCalled(t, "Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)"))
	})

	t.Run("error in rep", func(t *testing.T) {
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
		require.ErrorIs(t, err, expectedErr)
		mockAuthorRep.AssertCalled(t, "Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)"))
	})

	t.Run("error in update function", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()
		updateReq := models.AuthorUpdateReq{}

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)")).
			Return(errors.New("validation error"))

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Update(ctx, authorID, updateReq)
		require.Error(t, err)
		mockAuthorRep.AssertCalled(t, "Update", ctx, authorID, mock.AnythingOfType("func(*models.Author) (*models.Author, error)"))
	})

}

func TestAuthorServ_Delete(t *testing.T) {
	t.Run("success delete author", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(false, nil)
		mockAuthorRep.On("Delete", ctx, authorID).Return(nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)
		require.Nil(t, err)
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertCalled(t, "Delete", ctx, authorID)
	})

	t.Run("error author has linked artworks", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(true, nil)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)
		require.ErrorIs(t, err, authorserv.ErrHasLinkedArtworks)
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertNotCalled(t, "Delete", ctx, authorID)
	})

	t.Run("error checking artworks", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()
		expectedErr := errors.New("check artworks error")

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(false, expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "authorServ.Delete")
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertNotCalled(t, "Delete", ctx, authorID)
	})

	t.Run("error in delete", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()
		expectedErr := errors.New("delete error")

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockAuthorRep.On("HasArtworks", ctx, authorID).Return(false, nil)
		mockAuthorRep.On("Delete", ctx, authorID).Return(expectedErr)

		authorServ := authorserv.NewAuthorServ(mockAuthorRep)
		// ACT
		err := authorServ.Delete(ctx, authorID)
		require.ErrorIs(t, err, expectedErr)
		mockAuthorRep.AssertCalled(t, "HasArtworks", ctx, authorID)
		mockAuthorRep.AssertCalled(t, "Delete", ctx, authorID)
	})
}
