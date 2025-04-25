package searcher

import (
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep/mockartworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep/mockeventrep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	)
	return &event
}

func TestSearcher_GetAllArtworks(t *testing.T) {
	tests := []struct {
		name           string
		mockArtworks   []*models.Artwork
		expectedLength int
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockArt.On("GetAll").Return(tt.mockArtworks)

			result := service.GetAllArtworks()
			assert.Equal(t, tt.expectedLength, len(result))
			mockArt.AssertExpectations(t)
		})
	}
}

func TestSearcher_FilterArtworkByTitle(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		mockArtworks  []*models.Artwork
		expectedCount int
	}{
		{
			name:          "exact match",
			title:         "Test Artwork",
			mockArtworks:  []*models.Artwork{createTestArtwork()},
			expectedCount: 1,
		},
		{
			name:          "no match",
			title:         "Nonexistent",
			mockArtworks:  []*models.Artwork{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockArt.On("GetByTitle", tt.title).Return(tt.mockArtworks)

			result := service.FilterArtworkByTitle(tt.title)
			assert.Equal(t, tt.expectedCount, len(result))
			mockArt.AssertExpectations(t)
		})
	}
}

func TestSearcher_FilterArtworkByAuthor(t *testing.T) {
	author := createTestAuthor()
	tests := []struct {
		name          string
		author        *models.Author
		mockArtworks  []*models.Artwork
		expectedCount int
	}{
		{
			name:          "author with artworks",
			author:        author,
			mockArtworks:  []*models.Artwork{createTestArtwork()},
			expectedCount: 1,
		},
		{
			name:          "author without artworks",
			author:        author,
			mockArtworks:  []*models.Artwork{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockArt.On("GetByAuthor", tt.author).Return(tt.mockArtworks)

			result := service.FilterArtworkByAuthor(tt.author)
			assert.Equal(t, tt.expectedCount, len(result))
			mockArt.AssertExpectations(t)
		})
	}
}

func TestSearcher_FilterArtworkByCreationTime(t *testing.T) {
	artwork := createTestArtwork()
	tests := []struct {
		name          string
		yearBeg       int
		yearEnd       int
		mockArtworks  []*models.Artwork
		expectedCount int
	}{
		{
			name:          "within range",
			yearBeg:       1940,
			yearEnd:       1960,
			mockArtworks:  []*models.Artwork{artwork},
			expectedCount: 1,
		},
		{
			name:          "outside range",
			yearBeg:       1970,
			yearEnd:       1980,
			mockArtworks:  []*models.Artwork{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockArt.On("GetByCreationTime", tt.yearBeg, tt.yearEnd).Return(tt.mockArtworks)

			result := service.FilterArtworkByCreationTime(tt.yearBeg, tt.yearEnd)
			assert.Equal(t, tt.expectedCount, len(result))
			mockArt.AssertExpectations(t)
		})
	}
}

func TestSearcher_FilterArtworkByEvent(t *testing.T) {
	event := createTestEvent()
	tests := []struct {
		name          string
		event         models.Event
		mockArtworks  []*models.Artwork
		expectedCount int
	}{
		{
			name:          "event with artworks",
			event:         *event,
			mockArtworks:  []*models.Artwork{createTestArtwork()},
			expectedCount: 1,
		},
		{
			name:          "event without artworks",
			event:         *event,
			mockArtworks:  []*models.Artwork{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockArt.On("GetByEvent", tt.event).Return(tt.mockArtworks)

			result := service.FilterArtworkByEvent(tt.event)
			assert.Equal(t, tt.expectedCount, len(result))
			mockArt.AssertExpectations(t)
		})
	}
}

func TestSearcher_GetAllEvents(t *testing.T) {
	tests := []struct {
		name          string
		mockEvents    []*models.Event
		expectedCount int
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockEvent.On("GetAll").Return(tt.mockEvents)

			result := service.GetAllEvents()
			assert.Equal(t, tt.expectedCount, len(result))
			mockEvent.AssertExpectations(t)
		})
	}
}

func TestSearcher_FilterEventsByDate(t *testing.T) {
	event := createTestEvent()
	tests := []struct {
		name          string
		dateBeg       time.Time
		dateEnd       time.Time
		mockEvents    []*models.Event
		expectedCount int
	}{
		{
			name:          "within date range",
			dateBeg:       event.GetDateBegin().Add(-24 * time.Hour),
			dateEnd:       event.GetDateEnd().Add(24 * time.Hour),
			mockEvents:    []*models.Event{event},
			expectedCount: 1,
		},
		{
			name:          "outside date range",
			dateBeg:       event.GetDateBegin().Add(-48 * time.Hour),
			dateEnd:       event.GetDateBegin().Add(-24 * time.Hour),
			mockEvents:    []*models.Event{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockEvent.On("GetByDate", tt.dateBeg, tt.dateEnd).Return(tt.mockEvents)

			result := service.FilterEventsByDate(tt.dateBeg, tt.dateEnd)
			assert.Equal(t, tt.expectedCount, len(result))
			mockEvent.AssertExpectations(t)
		})
	}
}

func TestSearcher_GetEventOfArtworkOnDate(t *testing.T) {
	event := createTestEvent()
	artwork := createTestArtwork()
	dateBeg := event.GetDateBegin().Add(-24 * time.Hour)
	dateEnd := event.GetDateEnd().Add(24 * time.Hour)

	tests := []struct {
		name           string
		artwork        *models.Artwork
		dateBeg        time.Time
		dateEnd        time.Time
		mockEvent      *models.Event
		mockError      error
		expectedError  error
		expectedResult *models.Event
	}{
		{
			name:           "event found",
			artwork:        artwork,
			dateBeg:        dateBeg,
			dateEnd:        dateEnd,
			mockEvent:      event,
			mockError:      nil,
			expectedError:  nil,
			expectedResult: event,
		},
		{
			name:           "event not found",
			artwork:        artwork,
			dateBeg:        event.GetDateBegin().Add(-48 * time.Hour),
			dateEnd:        event.GetDateBegin().Add(-24 * time.Hour),
			mockEvent:      nil,
			mockError:      eventrep.ErrEventNotFound,
			expectedError:  eventrep.ErrEventNotFound,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArt := new(mockartworkrep.MockArtworkRep)
			mockEvent := new(mockeventrep.MockEventRep)
			service := NewSearcher(mockArt, mockEvent)

			mockEvent.On("GetEventOfArtworkOnDate", tt.artwork, tt.dateBeg, tt.dateEnd).Return(tt.mockEvent, tt.mockError)

			result, err := service.GetEventOfArtworkOnDate(tt.artwork, tt.dateBeg, tt.dateEnd)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			mockEvent.AssertExpectations(t)
		})
	}
}
