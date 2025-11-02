package collectionserv_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/collectionserv"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type CollectionServiceSuite struct {
	suite.Suite
}

func TestCollectionService(t *testing.T) {
	suite.RunSuite(t, new(CollectionServiceSuite))
}

func (s *CollectionServiceSuite) TestCollectionServ_GetAll(t provider.T) {
	collectionCreator := testobj.NewCollectionMother()

	t.WithNewStep("success return 2 collections", func(sCtx provider.StepCtx) {
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
		sCtx.Require().NoError(err)
		sCtx.Require().True(len(collections) == len(resCols))
		for i := range len(resCols) {
			sCtx.Require().True(collections[i].Equal(resCols[i]))
		}
		mockColRes.AssertCalled(t, "GetAll", ctx)
	})

	t.WithNewStep("success return 0 collections", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collections := []*models.Collection{}
		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("GetAll", ctx).Return(collections, nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		resCols, err := collectionServ.GetAll(ctx)
		sCtx.Require().NoError(err)
		sCtx.Require().True(len(resCols) == 0)
		mockColRep.AssertCalled(t, "GetAll", ctx)
	})

	t.WithNewStep("error in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collections := []*models.Collection{}
		expectedErr := errors.New("collectionRep error")
		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("GetAll", ctx).Return(collections, expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		_, err := collectionServ.GetAll(ctx)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockColRep.AssertCalled(t, "GetAll", ctx)
	})
}

func (s *CollectionServiceSuite) TestCollectionServ_Add(t provider.T) {
	collectionCreator := testobj.NewCollectionMother()

	t.WithNewStep("success add collection", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collection := collectionCreator.CollectionP(uuid.New())

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Add", ctx, collection).Return(nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Add(ctx, collection)
		sCtx.Require().NoError(err)
		mockColRep.AssertCalled(t, "Add", ctx, collection)
	})

	t.WithNewStep("error in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collection := collectionCreator.CollectionP(uuid.New())
		expectedErr := errors.New("database error")

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Add", ctx, collection).Return(expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Add(ctx, collection)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockColRep.AssertCalled(t, "Add", ctx, collection)
	})
}

func (s *CollectionServiceSuite) TestCollectionServ_Update(t provider.T) {
	collectionCreator := testobj.NewCollectionMother()

	t.WithNewStep("success update collection", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collectionID := uuid.New()
		updateReq := collectionCreator.CollectionUpdateReq()

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)")).Return(nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Update(ctx, collectionID, updateReq)
		sCtx.Require().NoError(err)
		mockColRep.AssertCalled(t, "Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)"))
	})

	t.WithNewStep("error in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collectionID := uuid.New()
		updateReq := collectionCreator.CollectionUpdateReq()
		expectedErr := errors.New("update error")

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)")).Return(expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Update(ctx, collectionID, updateReq)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockColRep.AssertCalled(t, "Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)"))
	})

	t.WithNewStep("error in update function", func(sCtx provider.StepCtx) {
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
		sCtx.Require().Error(err)
		mockColRep.AssertCalled(t, "Update", ctx, collectionID, mock.AnythingOfType("func(*models.Collection) (*models.Collection, error)"))
	})

}

func (s *CollectionServiceSuite) TestCollectionServ_Delete(t provider.T) {

	t.WithNewStep("success delete collection", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collectionID := uuid.New()

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Delete", ctx, collectionID).Return(nil)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Delete(ctx, collectionID)
		sCtx.Require().NoError(err)
		mockColRep.AssertCalled(t, "Delete", ctx, collectionID)
	})

	t.WithNewStep("error in delete", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		collectionID := uuid.New()
		expectedErr := errors.New("delete error")

		mockColRep := new(collectionrep.MockCollectionRep)
		mockColRep.On("Delete", ctx, collectionID).Return(expectedErr)

		collectionServ := collectionserv.NewCollectionServ(mockColRep)
		// ACT
		err := collectionServ.Delete(ctx, collectionID)
		sCtx.Require().ErrorIs(err, expectedErr)
		mockColRep.AssertCalled(t, "Delete", ctx, collectionID)
	})
}
