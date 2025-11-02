package searcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/repository/artworkrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/eventrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/searcher"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type SearcherServiceSuite struct {
	suite.Suite
}

func TestSearcherService(t *testing.T) {
	suite.RunSuite(t, new(SearcherServiceSuite))
}

func (s *SearcherServiceSuite) TestSearcher_GetAllArtworks(t provider.T) {
	// eventCreator := testobj.NewEventMother()
	artworkCreator := testobj.NewArtworkMother()

	t.WithNewStep("success return 2 artowrks", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		artworks := []*models.Artwork{
			artworkCreator.ArtworkP(uuid.New()),
			artworkCreator.ArtworkP(uuid.New()),
		}
		artworkFilter := artworkCreator.ArtworkFilterEmpty()
		artworkSortOps := artworkCreator.ArtworkSortOps()

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockArtworkRep.On(
			"GetAllArtworks", ctx,
			artworkFilter,
			artworkSortOps,
		).Return(artworks, nil)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		resArtworks, err := searcherServ.GetAllArtworks(ctx, artworkFilter, artworkSortOps)

		sCtx.Require().NoError(err)
		sCtx.Assert().True(len(artworks) == len(resArtworks))
		for i := range len(resArtworks) {
			sCtx.Assert().True(artworks[i].Equal(resArtworks[i]))
		}
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx,
			artworkFilter,
			artworkSortOps)
	})

	t.WithNewStep("success return 0 artworks", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		artworks := []*models.Artwork{}
		artworkFilter := artworkCreator.ArtworkFilterEmpty()
		artworkSortOps := artworkCreator.ArtworkSortOps()

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockArtworkRep.On(
			"GetAllArtworks", ctx,
			artworkFilter,
			artworkSortOps,
		).Return(artworks, nil)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		resArtworks, err := searcherServ.GetAllArtworks(ctx, artworkFilter, artworkSortOps)

		sCtx.Require().NoError(err)
		sCtx.Assert().True(len(resArtworks) == 0)
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx, artworkFilter, artworkSortOps)
	})

	t.WithNewStep("error in artwork rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		artworks := []*models.Artwork{}
		artworkFilter := artworkCreator.ArtworkFilterEmpty()
		artworkSortOps := artworkCreator.ArtworkSortOps()
		expectedErr := errors.New("artwork rep error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockArtworkRep.On(
			"GetAllArtworks", ctx,
			artworkFilter,
			artworkSortOps,
		).Return(artworks, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetAllArtworks(ctx, artworkFilter, artworkSortOps)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx, artworkFilter, artworkSortOps)
	})

}

func (s *SearcherServiceSuite) TestSearcher_GetAllEvents(t provider.T) {
	eventCreator := testobj.NewEventMother()

	t.WithNewStep("success return 2 events", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{
			eventCreator.EventP(uuid.New()),
			eventCreator.EventP(uuid.New()),
		}
		eventFilter := &jsonreqresp.EventFilter{
			DateBegin: time.Now(),
			DateEnd:   time.Now().AddDate(0, 0, 7),
		}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On("GetAll", ctx, eventFilter).Return(events, nil)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		resEvents, err := searcherServ.GetAllEvents(ctx, eventFilter)

		sCtx.Require().NoError(err)
		sCtx.Assert().True(len(events) == len(resEvents))
		for i := range len(resEvents) {
			sCtx.Assert().True(events[i].Equal(resEvents[i]))
		}
		mockEventRep.AssertCalled(t, "GetAll", ctx, eventFilter)
	})

	t.WithNewStep("error invalid date filter", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventFilter := &jsonreqresp.EventFilter{
			DateBegin: time.Now().AddDate(0, 0, 7), // Future date
			DateEnd:   time.Now(),                  // Past date
		}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetAllEvents(ctx, eventFilter)

		sCtx.Assert().Error(err)
		sCtx.Assert().ErrorIs(err, jsonreqresp.ErrEventFilterDate)
		mockEventRep.AssertNotCalled(t, "GetAll", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error in event rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{}
		eventFilter := &jsonreqresp.EventFilter{
			DateBegin: time.Now(),
			DateEnd:   time.Now().AddDate(0, 0, 7),
		}
		expectedErr := errors.New("event rep error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On("GetAll", ctx, eventFilter).Return(events, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetAllEvents(ctx, eventFilter)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetAll", ctx, eventFilter)
	})
}

func (s *SearcherServiceSuite) TestSearcher_GetEvent(t provider.T) {
	eventCreator := testobj.NewEventMother()

	t.WithNewStep("success get event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		event := eventCreator.EventP(eventID)

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On("GetByID", ctx, eventID).Return(event, nil)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		resEvent, err := searcherServ.GetEvent(ctx, eventID)

		sCtx.Require().NoError(err)
		sCtx.Assert().True(event.Equal(resEvent))
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
	})

	t.WithNewStep("error event not found", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("event not found")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On("GetByID", ctx, eventID).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetEvent(ctx, eventID)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
	})
}

func (s *SearcherServiceSuite) TestSearcher_GetArtworksFromEvent(t provider.T) {
	artworkCreator := testobj.NewArtworkMother()
	// eventCreator := testobj.NewEventMother()

	t.WithNewStep("success get artworks from event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkIDs := uuid.UUIDs{uuid.New(), uuid.New()}
		artworks := []*models.Artwork{
			artworkCreator.ArtworkP(artworkIDs[0]),
			artworkCreator.ArtworkP(artworkIDs[1]),
		}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(artworkIDs, nil)
		mockArtworkRep.On("GetByID", ctx, artworkIDs[0]).Return(artworks[0], nil)
		mockArtworkRep.On("GetByID", ctx, artworkIDs[1]).Return(artworks[1], nil)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		resArtworks, err := searcherServ.GetArtworksFromEvent(ctx, eventID)

		sCtx.Require().NoError(err)
		sCtx.Assert().True(len(artworks) == len(resArtworks))
		for i := range len(resArtworks) {
			sCtx.Assert().True(artworks[i].Equal(resArtworks[i]))
		}
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[1])
	})

	t.WithNewStep("error getting artwork IDs", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("get artwork IDs error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetArtworksFromEvent(ctx, eventID)

		sCtx.Assert().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error getting artwork by ID", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkIDs := uuid.UUIDs{uuid.New()}
		expectedErr := errors.New("artwork not found")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(artworkIDs, nil)
		mockArtworkRep.On("GetByID", ctx, artworkIDs[0]).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetArtworksFromEvent(ctx, eventID)

		sCtx.Assert().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
	})

}

func (s *SearcherServiceSuite) TestSearcher_GetCollectionsStat(t provider.T) {
	eventCreator := testobj.NewEventMother()

	t.WithNewStep("success get collections stat", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		statCols := []*models.StatCollections{
			eventCreator.StatCollectionsP(),
			eventCreator.StatCollectionsP(),
		}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		event := eventCreator.EventP(eventID)
		mockEventRep.On("GetByID", ctx, eventID).Return(event, nil)
		mockEventRep.On("GetCollectionsStat", ctx, eventID).Return(statCols, nil)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		resStatCols, err := searcherServ.GetCollectionsStat(ctx, eventID)

		sCtx.Require().NoError(err)
		sCtx.Assert().True(len(statCols) == len(resStatCols))
		for i := range len(resStatCols) {
			sCtx.Assert().True(statCols[i].Equals(resStatCols[i]))
		}
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetCollectionsStat", ctx, eventID)
	})

	t.WithNewStep("error event not found", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("event not found")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetByID", ctx, eventID).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetCollectionsStat(ctx, eventID)

		sCtx.Assert().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertNotCalled(t, "GetCollectionsStat", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error getting collections stat", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("stat error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		event := eventCreator.EventP(eventID)
		mockEventRep.On("GetByID", ctx, eventID).Return(event, nil)
		mockEventRep.On("GetCollectionsStat", ctx, eventID).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetCollectionsStat(ctx, eventID)

		sCtx.Assert().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetCollectionsStat", ctx, eventID)
	})
}
