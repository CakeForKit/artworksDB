package searcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/searcher"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSearcher_GetAllArtworks(t *testing.T) {
	// eventCreator := testobj.NewEventMother()
	artworkCreator := testobj.NewArtworkMother()

	t.Run("success return 2 artowrks", func(t *testing.T) {
		ctx := context.Background()
		artworks := []*models.Artwork{
			artworkCreator.ArtworkP(uuid.New()),
			artworkCreator.ArtworkP(uuid.New()),
		}
		artworkFilter := artworkCreator.ArtworkFilter()
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

		require.Nil(t, err)
		require.True(t, len(artworks) == len(resArtworks))
		for i := range len(resArtworks) {
			require.True(t, artworks[i].Equals(resArtworks[i]))
		}
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx,
			artworkFilter,
			artworkSortOps)
	})

	t.Run("success return 0 artworks", func(t *testing.T) {
		ctx := context.Background()
		artworks := []*models.Artwork{}
		artworkFilter := artworkCreator.ArtworkFilter()
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

		require.Nil(t, err)
		require.True(t, len(resArtworks) == 0)
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx, artworkFilter, artworkSortOps)
	})

	t.Run("error in artwork rep", func(t *testing.T) {
		ctx := context.Background()
		artworks := []*models.Artwork{}
		artworkFilter := artworkCreator.ArtworkFilter()
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

		require.ErrorIs(t, err, expectedErr)
		mockArtworkRep.AssertCalled(t, "GetAllArtworks", ctx, artworkFilter, artworkSortOps)
	})

}

func TestSearcher_GetAllEvents(t *testing.T) {
	eventCreator := testobj.NewEventMother()

	t.Run("success return 2 events", func(t *testing.T) {
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

		require.Nil(t, err)
		require.True(t, len(events) == len(resEvents))
		for i := range len(resEvents) {
			require.True(t, events[i].Equals(resEvents[i]))
		}
		mockEventRep.AssertCalled(t, "GetAll", ctx, eventFilter)
	})

	t.Run("error invalid date filter", func(t *testing.T) {
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

		require.Error(t, err)
		require.ErrorIs(t, err, jsonreqresp.ErrEventFilterDate)
		mockEventRep.AssertNotCalled(t, "GetAll", mock.Anything, mock.Anything)
	})

	t.Run("error in event rep", func(t *testing.T) {
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

		require.ErrorIs(t, err, expectedErr)
		mockEventRep.AssertCalled(t, "GetAll", ctx, eventFilter)
	})
}

func TestSearcher_GetEvent(t *testing.T) {
	eventCreator := testobj.NewEventMother()

	t.Run("success get event", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()
		event := eventCreator.EventP(eventID)

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On("GetByID", ctx, eventID).Return(event, nil)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		resEvent, err := searcherServ.GetEvent(ctx, eventID)

		require.Nil(t, err)
		require.True(t, event.Equals(resEvent))
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
	})

	t.Run("error event not found", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("event not found")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On("GetByID", ctx, eventID).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetEvent(ctx, eventID)

		require.ErrorIs(t, err, expectedErr)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
	})
}

func TestSearcher_GetArtworksFromEvent(t *testing.T) {
	artworkCreator := testobj.NewArtworkMother()
	// eventCreator := testobj.NewEventMother()

	t.Run("success get artworks from event", func(t *testing.T) {
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

		require.Nil(t, err)
		require.True(t, len(artworks) == len(resArtworks))
		for i := range len(resArtworks) {
			require.True(t, artworks[i].Equals(resArtworks[i]))
		}
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[1])
	})

	t.Run("error getting artwork IDs", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("get artwork IDs error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetArtworksFromEvent(ctx, eventID)

		require.Error(t, err)
		require.Contains(t, err.Error(), "searcher.GetArtworkFromEvent")
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("error getting artwork by ID", func(t *testing.T) {
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

		require.Error(t, err)
		require.Contains(t, err.Error(), "searcher.GetArtworkFromEvent")
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
	})

}

func TestSearcher_GetCollectionsStat(t *testing.T) {
	eventCreator := testobj.NewEventMother()

	t.Run("success get collections stat", func(t *testing.T) {
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

		require.Nil(t, err)
		require.True(t, len(statCols) == len(resStatCols))
		for i := range len(resStatCols) {
			require.True(t, statCols[i].Equals(resStatCols[i]))
		}
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetCollectionsStat", ctx, eventID)
	})

	t.Run("error event not found", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("event not found")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetByID", ctx, eventID).Return(nil, expectedErr)

		searcherServ := searcher.NewSearcher(mockArtworkRep, mockEventRep)

		// ACT
		_, err := searcherServ.GetCollectionsStat(ctx, eventID)

		require.Error(t, err)
		require.Contains(t, err.Error(), "searcher.GetCollectionsStat")
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertNotCalled(t, "GetCollectionsStat", mock.Anything, mock.Anything)
	})

	t.Run("error getting collections stat", func(t *testing.T) {
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

		require.Error(t, err)
		require.Contains(t, err.Error(), "searcher.GetCollectionsStat")
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetCollectionsStat", ctx, eventID)
	})
}

/*
func createTestAuthor() *models.Author {
	author, _ := models.NewAuthor(uuid.New(), "Test Author", 1900, 2000)
	return &author
}

func createTestCollection() *models.Collection {
	collection, _ := models.NewCollection(uuid.New(), "Test Collection")
	return &collection
}

func createTestArtwork() *models.Artwork {
	author := createTestAuthor()
	collection := createTestCollection()
	artwork, _ := models.NewArtwork(
		uuid.New(),
		"Test Artwork",
		"oil on canvas",
		"canvas",
		"100x100 cm",
		1950,
		author,
		collection,
	)
	return &artwork
}

func createTestEvent() *models.Event {
	event, _ := models.NewEvent(
		uuid.New(),
		"Test Event",
		time.Now(),
		time.Now().Add(24*time.Hour),
		"Test Address",
		true,
		uuid.New(),
		100,
		true,
		make(uuid.UUIDs, 0),
	)
	return &event
}

func TestSearcher_GetAllArtworks(t *testing.T) {
	ctx := context.Background()
	filter := &jsonreqresp.ArtworkFilter{}
	sort := &jsonreqresp.ArtworkSortOps{}

	tests := []struct {
		name           string
		mockArtworks   []*models.Artwork
		mockError      error
		expectedLength int
		expectedError  error
	}{
		{
			name:           "single artwork",
			mockArtworks:   []*models.Artwork{createTestArtwork()},
			expectedLength: 1,
		},
		{
			name:           "multiple artworks",
			mockArtworks:   []*models.Artwork{createTestArtwork(), createTestArtwork()},
			expectedLength: 2,
		},
		{
			name:          "repository error",
			mockError:     artworkrep.ErrArtworkNotFound,
			expectedError: artworkrep.ErrArtworkNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := &artworkrep.MockArtworkRep{}
			mockEvent := &eventrep.MockEventRep{}
			service := searcher.NewSearcher(mockArt, mockEvent)

			mockArt.On("GetAllArtworks", ctx, filter, sort).Return(tt.mockArtworks, tt.mockError)

			result, err := service.GetAllArtworks(ctx, filter, sort)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedLength, len(result))
			}

			mockArt.AssertExpectations(t)
		})
	}
}

func TestSearcher_GetAllEvents(t *testing.T) {
	ctx := context.Background()
	filter := &jsonreqresp.EventFilter{
		DateBegin: time.Now(),
		DateEnd:   time.Now().Add(24 * time.Hour),
	}

	tests := []struct {
		name          string
		mockEvents    []*models.Event
		mockError     error
		expectedCount int
		expectedError error
	}{
		{
			name:          "single event",
			mockEvents:    []*models.Event{createTestEvent()},
			expectedCount: 1,
		},
		{
			name:          "multiple events",
			mockEvents:    []*models.Event{createTestEvent(), createTestEvent()},
			expectedCount: 2,
		},
		{
			name:          "repository error",
			mockError:     eventrep.ErrEventNotFound,
			expectedError: eventrep.ErrEventNotFound,
		},
		{
			name:          "invalid date range",
			mockError:     jsonreqresp.ErrEventFilterDate,
			expectedError: jsonreqresp.ErrEventFilterDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := &artworkrep.MockArtworkRep{}
			mockEvent := &eventrep.MockEventRep{}
			service := searcher.NewSearcher(mockArt, mockEvent)

			if tt.name == "invalid date range" {
				invalidFilter := &jsonreqresp.EventFilter{
					DateBegin: time.Now().Add(24 * time.Hour),
					DateEnd:   time.Now(),
				}
				result, err := service.GetAllEvents(ctx, invalidFilter)
				assert.ErrorIs(t, err, jsonreqresp.ErrEventFilterDate)
				assert.Nil(t, result)
				return
			}

			mockEvent.On("GetAll", ctx, filter).Return(tt.mockEvents, tt.mockError)

			result, err := service.GetAllEvents(ctx, filter)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result))
			}

			mockEvent.AssertExpectations(t)
		})
	}
}
*/
