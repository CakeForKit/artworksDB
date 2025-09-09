package collectionserv_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/collectionserv"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthorServ_GetAll(t *testing.T) {
	collectionCreator := testobj.NewCollectionMother()

	t.Run("success return 2 collections", func(t *testing.T) {
		ctx := context.Background()
		collections := []*models.Collection{
			collectionCreator.CollectionP(uuid.New()),
			collectionCreator.CollectionP(uuid.New()),
		}
		mockColRes := new(collectionrep.MockCollectionRep)
		mockColRes.On("GetAll", ctx).Return(collections, nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRes)
		// ACT
		resCols, err := collectionServ.GetAll(ctx)
		require.Nil(t, err)
		require.True(t, len(collections) == len(resCols))
		for i := range len(resCols) {
			require.True(t, collections[i].Equals(resCols[i]))
		}
		mockColRes.AssertCalled(t, "GetAll", ctx)
	})

	t.Run("success return 0 collections", func(t *testing.T) {
		ctx := context.Background()
		collections := []*models.Collection{}
		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("GetAll", ctx).Return(collections, nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		resCols, err := collectionServ.GetAll(ctx)
		require.Nil(t, err)
		require.True(t, len(resCols) == 0)
		mockColRep.AssertCalled(t, "GetAll", ctx)
	})

	t.Run("error in rep", func(t *testing.T) {
		ctx := context.Background()
		collections := []*models.Collection{}
		expectedErr := errors.New("collectionRep error")
		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("GetAll", ctx).Return(collections, expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		_, err := collectionServ.GetAll(ctx)
		require.ErrorIs(t, err, expectedErr)
		mockColRep.AssertCalled(t, "GetAll", ctx)
	})
}

func TestCollectionServ_Add(t *testing.T) {
	collectionCreator := testobj.NewCollectionMother()

	t.Run("success add collection", func(t *testing.T) {
		ctx := context.Background()
		collection := collectionCreator.CollectionP(uuid.New())

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Add", ctx, collection).Return(nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Add(ctx, collection)
		require.Nil(t, err)
		mockColRep.AssertCalled(t, "Add", ctx, collection)
	})

	t.Run("error in rep", func(t *testing.T) {
		ctx := context.Background()
		collection := collectionCreator.CollectionP(uuid.New())
		expectedErr := errors.New("database error")

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Add", ctx, collection).Return(expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Add(ctx, collection)
		require.ErrorIs(t, err, expectedErr)
		mockColRep.AssertCalled(t, "Add", ctx, collection)
	})
}

func TestCollectionServ_Update(t *testing.T) {
	collectionCreator := testobj.NewCollectionMother()

	t.Run("success update collection", func(t *testing.T) {
		ctx := context.Background()
		collectionID := uuid.New()
		updateReq := collectionCreator.CollectionUpdateReq()

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)")).Return(nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Update(ctx, collectionID, updateReq)
		require.Nil(t, err)
		mockColRep.AssertCalled(t, "Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)"))
	})

	t.Run("error in rep", func(t *testing.T) {
		ctx := context.Background()
		collectionID := uuid.New()
		updateReq := collectionCreator.CollectionUpdateReq()
		expectedErr := errors.New("update error")

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)")).Return(expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Update(ctx, collectionID, updateReq)
		require.ErrorIs(t, err, expectedErr)
		mockColRep.AssertCalled(t, "Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)"))
	})

	t.Run("error in update function", func(t *testing.T) {
		ctx := context.Background()
		collectionID := uuid.New()
		// Invalid update request that should cause error in Update method
		updateReq := models.CollectionUpdateReq{}

		mockColRep := new(collectionrep.MockCollectionRep)
		// Mock will call the actual update function which should fail
		mockColRep.On("Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)")).
			Return(errors.New("validation error"))

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Update(ctx, collectionID, updateReq)
		require.Error(t, err)
		mockColRep.AssertCalled(t, "Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)"))
	})

}

func TestCollectionServ_Delete(t *testing.T) {

	t.Run("success delete collection", func(t *testing.T) {
		ctx := context.Background()
		collectionID := uuid.New()

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Delete", ctx, collectionID).Return(nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Delete(ctx, collectionID)
		require.Nil(t, err)
		mockColRep.AssertCalled(t, "Delete", ctx, collectionID)
	})

	t.Run("error in delete", func(t *testing.T) {
		ctx := context.Background()
		collectionID := uuid.New()
		expectedErr := errors.New("delete error")

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Delete", ctx, collectionID).Return(expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Delete(ctx, collectionID)
		require.ErrorIs(t, err, expectedErr)
		mockColRep.AssertCalled(t, "Delete", ctx, collectionID)
	})
}
