package artworkserv_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthorServ_GetAll(t *testing.T) {
	artworkCretor := testobj.NewArtworkMother()

	t.Run("success return 2 authors", func(t *testing.T) {
		ctx := context.Background()
		artworks := []*models.Artwork{
			artworkCretor.ArtworkP(uuid.New()),
			artworkCretor.ArtworkP(uuid.New()),
		}
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockArtworkRep.On(
			"GetAllArtworks", ctx,
			mock.AnythingOfType("*jsonreqresp.ArtworkFilter"),
			mock.AnythingOfType("*jsonreqresp.ArtworkSortOps"),
		).Return(artworks, nil)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		resArtworks, err := artworkService.GetAll(ctx)

		require.Nil(t, err)
		require.True(t, len(artworks) == len(resArtworks))
		for i := range len(resArtworks) {
			require.True(t, artworks[i].Equals(resArtworks[i]))
		}
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx,
			mock.AnythingOfType("*jsonreqresp.ArtworkFilter"),
			mock.AnythingOfType("*jsonreqresp.ArtworkSortOps"))
	})

	t.Run("success return 0 artworks", func(t *testing.T) {
		ctx := context.Background()
		artworks := []*models.Artwork{}
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockArtworkRep.On(
			"GetAllArtworks", ctx,
			mock.AnythingOfType("*jsonreqresp.ArtworkFilter"),
			mock.AnythingOfType("*jsonreqresp.ArtworkSortOps"),
		).Return(artworks, nil)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		resArtworks, err := artworkService.GetAll(ctx)

		require.Nil(t, err)
		require.True(t, len(resArtworks) == 0)
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx,
			mock.AnythingOfType("*jsonreqresp.ArtworkFilter"),
			mock.AnythingOfType("*jsonreqresp.ArtworkSortOps"))
	})

	t.Run("error in rep", func(t *testing.T) {
		ctx := context.Background()
		artworks := []*models.Artwork{}
		expectedErr := errors.New("artworkRep error")
		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockArtworkRep.On(
			"GetAllArtworks", ctx,
			mock.AnythingOfType("*jsonreqresp.ArtworkFilter"),
			mock.AnythingOfType("*jsonreqresp.ArtworkSortOps"),
		).Return(artworks, expectedErr)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		_, err := artworkService.GetAll(ctx)

		require.ErrorIs(t, err, expectedErr)
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx,
			mock.AnythingOfType("*jsonreqresp.ArtworkFilter"),
			mock.AnythingOfType("*jsonreqresp.ArtworkSortOps"))
	})
}

func TestArtworkService_Add(t *testing.T) {
	artworkCreator := testobj.NewArtworkMother()
	authorCreator := testobj.NewAuthorMother()
	collectionCreator := testobj.NewCollectionMother()

	t.Run("success add artwork", func(t *testing.T) {
		ctx := context.Background()
		author := authorCreator.AuthorP(uuid.New())
		collection := collectionCreator.CollectionP(uuid.New())
		artworkReq := artworkCreator.AddArtworkRequest(author.GetBirthYear()+1, author.GetID(), collection.GetID())

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		mockAuthorRep.On("GetByID", ctx, author.GetID()).Return(author, nil)
		mockCollectionRep.On("GetByID", ctx, collection.GetID()).Return(collection, nil)
		mockArtworkRep.On("Add", ctx, mock.AnythingOfType("*models.Artwork")).Return(nil)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Add(ctx, artworkReq)

		require.Nil(t, err)
		mockAuthorRep.AssertCalled(t, "GetByID", ctx, author.GetID())
		mockCollectionRep.AssertCalled(t, "GetByID", ctx, collection.GetID())
		mockArtworkRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Artwork"))
	})

	t.Run("error author not found", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()
		collection := collectionCreator.CollectionP(uuid.New())
		artworkReq := artworkCreator.AddArtworkRequest(2000, authorID, collection.GetID())

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		expectedErr := errors.New("author not found")
		mockAuthorRep.On("GetByID", ctx, authorID).Return(nil, expectedErr)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Add(ctx, artworkReq)

		require.Error(t, err)
		require.Contains(t, err.Error(), "artworkService.Add")
		mockAuthorRep.AssertCalled(t, "GetByID", ctx, authorID)
		mockCollectionRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		mockArtworkRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("error collection not found", func(t *testing.T) {
		ctx := context.Background()
		author := authorCreator.AuthorP(uuid.New())
		collectionID := uuid.New()
		artworkReq := artworkCreator.AddArtworkRequest(author.GetBirthYear()+1, author.GetID(), collectionID)

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		mockAuthorRep.On("GetByID", ctx, author.GetID()).Return(author, nil)
		expectedErr := errors.New("collection not found")
		mockCollectionRep.On("GetByID", ctx, collectionID).Return(nil, expectedErr)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Add(ctx, artworkReq)

		require.Error(t, err)
		require.Contains(t, err.Error(), "artworkService.Add")
		mockAuthorRep.AssertCalled(t, "GetByID", ctx, author.GetID())
		mockCollectionRep.AssertCalled(t, "GetByID", ctx, collectionID)
		mockArtworkRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("error in artwork rep", func(t *testing.T) {
		ctx := context.Background()
		author := authorCreator.AuthorP(uuid.New())
		collection := collectionCreator.CollectionP(uuid.New())
		artworkReq := artworkCreator.AddArtworkRequest(author.GetBirthYear()+1, author.GetID(), collection.GetID())

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		mockAuthorRep.On("GetByID", ctx, author.GetID()).Return(author, nil)
		mockCollectionRep.On("GetByID", ctx, collection.GetID()).Return(collection, nil)
		expectedErr := errors.New("database error")
		mockArtworkRep.On("Add", ctx, mock.AnythingOfType("*models.Artwork")).Return(expectedErr)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Add(ctx, artworkReq)

		require.Error(t, err)
		require.Contains(t, err.Error(), "artworkService.Add")
		mockAuthorRep.AssertCalled(t, "GetByID", ctx, author.GetID())
		mockCollectionRep.AssertCalled(t, "GetByID", ctx, collection.GetID())
		mockArtworkRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Artwork"))
	})

}

func TestArtworkService_Delete(t *testing.T) {
	t.Run("success delete artwork", func(t *testing.T) {
		ctx := context.Background()
		artworkID := uuid.New()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		mockArtworkRep.On("Delete", ctx, artworkID).Return(nil)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Delete(ctx, artworkID)

		require.Nil(t, err)
		mockArtworkRep.AssertCalled(t, "Delete", ctx, artworkID)
	})

	t.Run("error in delete", func(t *testing.T) {
		ctx := context.Background()
		artworkID := uuid.New()

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		expectedErr := errors.New("delete error")
		mockArtworkRep.On("Delete", ctx, artworkID).Return(expectedErr)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Delete(ctx, artworkID)

		require.ErrorIs(t, err, expectedErr)
		mockArtworkRep.AssertCalled(t, "Delete", ctx, artworkID)
	})
}

func TestArtworkService_Update(t *testing.T) {
	authorCreator := testobj.NewAuthorMother()
	collectionCreator := testobj.NewCollectionMother()

	t.Run("success update artwork", func(t *testing.T) {
		ctx := context.Background()
		artworkID := uuid.New()
		author := authorCreator.AuthorP(uuid.New())
		collection := collectionCreator.CollectionP(uuid.New())
		updateFields := jsonreqresp.ArtworkUpdate{
			Title:        "Updated Title",
			AuthorID:     author.GetID().String(),
			CollectionID: collection.GetID().String(),
		}

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		mockAuthorRep.On("GetByID", ctx, author.GetID()).Return(author, nil)
		mockCollectionRep.On("GetByID", ctx, collection.GetID()).Return(collection, nil)
		mockArtworkRep.On("Update", ctx, artworkID, mock.AnythingOfType("func(*models.Artwork) (*models.Artwork, error)")).Return(nil)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Update(ctx, artworkID, updateFields)

		require.Nil(t, err)
		mockAuthorRep.AssertCalled(t, "GetByID", ctx, author.GetID())
		mockCollectionRep.AssertCalled(t, "GetByID", ctx, collection.GetID())
		mockArtworkRep.AssertCalled(t, "Update", ctx, artworkID, mock.AnythingOfType("func(*models.Artwork) (*models.Artwork, error)"))
	})

	t.Run("error author not found in update", func(t *testing.T) {
		ctx := context.Background()
		artworkID := uuid.New()
		authorID := uuid.New()
		updateFields := jsonreqresp.ArtworkUpdate{
			AuthorID: authorID.String(),
		}

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		expectedErr := errors.New("author not found")
		mockAuthorRep.On("GetByID", ctx, authorID).Return(nil, expectedErr)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Update(ctx, artworkID, updateFields)

		require.Error(t, err)
		require.Contains(t, err.Error(), "artworkService.Add")
		mockAuthorRep.AssertCalled(t, "GetByID", ctx, authorID)
		mockCollectionRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		mockArtworkRep.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("error in artwork rep update", func(t *testing.T) {
		ctx := context.Background()
		artworkID := uuid.New()
		author := authorCreator.AuthorP(uuid.New())
		collection := collectionCreator.CollectionP(uuid.New())
		updateFields := jsonreqresp.ArtworkUpdate{
			AuthorID:     author.GetID().String(),
			CollectionID: collection.GetID().String(),
		}

		mockAuthorRep := new(authorrep.MockAuthorRep)
		mockCollectionRep := new(collectionrep.MockCollectionRep)
		mockArtworkRep := new(artworkrep.MockArtworkRep)

		mockAuthorRep.On("GetByID", ctx, author.GetID()).Return(author, nil)
		mockCollectionRep.On("GetByID", ctx, collection.GetID()).Return(collection, nil)
		expectedErr := errors.New("update error")
		mockArtworkRep.On("Update", ctx, artworkID, mock.AnythingOfType("func(*models.Artwork) (*models.Artwork, error)")).Return(expectedErr)

		artworkService := artworkserv.NewArtworkService(mockArtworkRep, mockAuthorRep, mockCollectionRep)

		// ACT
		err := artworkService.Update(ctx, artworkID, updateFields)

		require.ErrorIs(t, err, expectedErr)
		mockAuthorRep.AssertCalled(t, "GetByID", ctx, author.GetID())
		mockCollectionRep.AssertCalled(t, "GetByID", ctx, collection.GetID())
		mockArtworkRep.AssertCalled(t, "Update", ctx, artworkID, mock.AnythingOfType("func(*models.Artwork) (*models.Artwork, error)"))
	})
}
