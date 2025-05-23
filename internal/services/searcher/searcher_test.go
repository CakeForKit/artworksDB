package searcher_test

import (
	"context"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/searcher"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
