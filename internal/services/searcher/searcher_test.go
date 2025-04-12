package searcher

import (
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep/mockartworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep/mockeventrep"
	"github.com/google/uuid"
	"github.com/stateio/testify/assert"
)

// func TestSearcher_GetAllArtworks(t *testing.T) {
// 	type testCase struct {
// 		name string
// 		res  []*models.Artwork
// 	}

// 	mockArtworkRep := new(mockartworkrep.MockArtworkRep)
// 	mockArtworkRep.On("GetAllArtworks").Return(&[]models.Artwork{})
// 	searcherServ :=
// 		searcher{artowrkRep: mockArtworkRep}

// 	var testCases []testCase = []testCase{
// 		{
// 			name: "Empty",
// 			res:  []*models.Artwork{},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			res := searcherServ.GetAllArtworks()
// 			if len(res) != 0 {
// 				t.Errorf("Not empty array??")
// 			}
// 		})
// 	}
// }

func createTestAuthor() models.Author {
	author, _ := models.NewAuthor("Test Author", 1900, 2000)
	return author
}

func createTestCollection() models.Collection {
	collection, _ := models.NewCollection("Test Collection")
	return collection
}

func createTestArtwork() models.Artwork {
	author := createTestAuthor()
	collection := createTestCollection()
	artwork, _ := models.NewArtwork(
		uuid.New(),
		"Test Artwork",
		1950,
		&author,
		&collection,
		"100x100",
		"oil",
		"painting",
	)
	return artwork
}

func createTestEvent() models.Event {
	artwork := createTestArtwork()
	event, _ := models.NewEvent(
		"Test Event",
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
		"Test Address",
		true,
		[]*models.Artwork{&artwork},
	)
	return event
}

func TestSearcher_GetAllArtworks(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	artwork := createTestArtwork()
	expectedArtworks := []*models.Artwork{&artwork}

	mockArt.On("GetAll").Return(expectedArtworks)

	t.Run("get all artworks", func(t *testing.T) {
		result := searcherServ.GetAllArtworks()
		assert.Equal(t, expectedArtworks, result)
		assert.Len(t, result, 1)
		assert.Equal(t, artwork.GetTitle(), result[0].GetTitle())
		mockArt.AssertExpectations(t)
	})
}

func TestSearcher_FilterArtworkByTitle(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	artwork := createTestArtwork()
	expectedArtworks := []*models.Artwork{&artwork}

	mockArt.On("GetByTitle", artwork.GetTitle()).Return(expectedArtworks)

	t.Run("filter by exact title", func(t *testing.T) {
		result := searcherServ.FilterArtworkByTitle(artwork.GetTitle())
		assert.Equal(t, expectedArtworks, result)
		assert.Len(t, result, 1)
		assert.Equal(t, artwork.GetID(), result[0].GetID())
		mockArt.AssertExpectations(t)
	})
}

func TestSearcher_FilterArtworkByAuthor(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	artwork := createTestArtwork()
	author := *artwork.GetAuthor()
	expectedArtworks := []*models.Artwork{&artwork}

	mockArt.On("GetByAuthor", &author).Return(expectedArtworks)

	t.Run("filter by author", func(t *testing.T) {
		result := searcherServ.FilterArtworkByAuthor(&author)
		assert.Equal(t, expectedArtworks, result)
		assert.Len(t, result, 1)
		assert.Equal(t, artwork.GetID(), result[0].GetID())
		mockArt.AssertExpectations(t)
	})
}

func TestSearcher_FilterArtworkByCreationTime(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	artwork := createTestArtwork()
	yearBeg := artwork.GetCreationYear() - 10
	yearEnd := artwork.GetCreationYear() + 10
	expectedArtworks := []*models.Artwork{&artwork}

	mockArt.On("GetByCreationTime", yearBeg, yearEnd).Return(expectedArtworks)

	t.Run("filter by creation time range", func(t *testing.T) {
		result := searcherServ.FilterArtworkByCreationTime(yearBeg, yearEnd)
		assert.Equal(t, expectedArtworks, result)
		assert.Len(t, result, 1)
		assert.Equal(t, artwork.GetCreationYear(), result[0].GetCreationYear())
		mockArt.AssertExpectations(t)
	})
}

func TestSearcher_FilterArtworkByEvent(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	event := createTestEvent()
	expectedArtworks := event.GetArtworks()

	mockArt.On("GetByEvent", event).Return(expectedArtworks)

	t.Run("filter by event", func(t *testing.T) {
		result := searcherServ.FilterArtworkByEvent(event)
		assert.Equal(t, expectedArtworks, result)
		assert.Len(t, result, 1)
		mockArt.AssertExpectations(t)
	})
}

func TestSearcher_GetAllEvents(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	event := createTestEvent()
	expectedEvents := []*models.Event{&event}

	mockEvent.On("GetAll").Return(expectedEvents)

	t.Run("get all events", func(t *testing.T) {
		result := searcherServ.GetAllEvents()
		assert.Equal(t, expectedEvents, result)
		assert.Len(t, result, 1)
		assert.Equal(t, event.GetTitle(), result[0].GetTitle())
		mockEvent.AssertExpectations(t)
	})
}

func TestSearcher_FilterEventsByDate(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	event := createTestEvent()
	dateBeg := event.GetDateBegin().Add(-24 * time.Hour)
	dateEnd := event.GetDateEnd().Add(24 * time.Hour)
	expectedEvents := []*models.Event{&event}

	mockEvent.On("GetByDate", dateBeg, dateEnd).Return(expectedEvents)

	t.Run("filter events by date range", func(t *testing.T) {
		result := searcherServ.FilterEventsByDate(dateBeg, dateEnd)
		assert.Equal(t, expectedEvents, result)
		assert.Len(t, result, 1)
		assert.True(t, event.GetDateBegin().After(dateBeg))
		assert.True(t, event.GetDateEnd().Before(dateEnd))
		mockEvent.AssertExpectations(t)
	})
}

func TestSearcher_GetEventOfArtworkOnDate(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	mockEvent := new(mockeventrep.MockEventRep)
	searcherServ := NewSearcher(mockArt, mockEvent)

	event := createTestEvent()
	artwork := *event.GetArtworks()[0]
	dateBeg := event.GetDateBegin().Add(-24 * time.Hour)
	dateEnd := event.GetDateEnd().Add(24 * time.Hour)

	mockEvent.On("GetEventOfArtworkOnDate", &artwork, dateBeg, dateEnd).Return(&event, nil)

	t.Run("get event for artwork on date", func(t *testing.T) {
		result, err := searcherServ.GetEventOfArtworkOnDate(&artwork, dateBeg, dateEnd)
		assert.NoError(t, err)
		assert.Equal(t, &event, result)
		assert.True(t, result.CheckArtwork(artwork.GetID()))
		mockEvent.AssertExpectations(t)
	})

	nonEventDateBeg := event.GetDateBegin().Add(-48 * time.Hour)
	nonEventDateEnd := event.GetDateBegin().Add(-24 * time.Hour)
	var enil *models.Event = nil
	mockEvent.On("GetEventOfArtworkOnDate", &artwork, nonEventDateBeg, nonEventDateEnd).Return(enil, eventrep.ErrEventNotFound)

	t.Run("event not found for artwork", func(t *testing.T) {
		result, err := searcherServ.GetEventOfArtworkOnDate(&artwork, nonEventDateBeg, nonEventDateEnd)
		assert.EqualError(t, err, eventrep.ErrEventNotFound.Error())
		assert.Nil(t, result)
		mockEvent.AssertExpectations(t)
	})
}

// func TestSearcher_GetByID(t *testing.T) {
// 	type testCase struct {
// 		name        string
// 		id          uuid.UUID
// 		aw          *models.Artwork
// 		expectedErr error
// 	}

// 	paw, err := createTestArtwork()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	mock := new(mockartworkrep.MockArtworkRep)
// 	var idExist uuid.UUID = paw.GetID()
// 	mock.On("GetByID", idExist).Return(paw, nil)
// 	idNOTExist := uuid.New()
// 	var awnill *models.Artwork = nil
// 	mock.On("GetByID", idNOTExist).Return(awnill, artworkrep.ErrArtworkNotFound)

// 	searcherServ := searcher{artowrkRep: mock}
// 	var testCases []testCase = []testCase{
// 		{
// 			name:        "id exist",
// 			id:          idExist,
// 			aw:          paw,
// 			expectedErr: nil,
// 		},
// 		{
// 			name:        "id NOT exist",
// 			id:          idNOTExist,
// 			expectedErr: artworkrep.ErrArtworkNotFound,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			awres, err := searcherServ.GetByID(tc.id)
// 			if err != tc.expectedErr {
// 				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
// 			} else if err == nil && awres != tc.aw {
// 				t.Errorf("Expected artwork %v \ngot %v", tc.aw, awres)
// 			}
// 		})
// 	}
// }

// func TestSearcher_AddArtwork(t *testing.T) {
// 	t.Run("successful artwork addition", func(t *testing.T) {
// 		// Arrange
// 		mockRepo := new(mockartworkrep.MockArtworkRep)
// 		searcher := &searcher{rep: mockRepo}
// 		var err error
// 		ptestArtwork, err := createTestArtwork()
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		mockRepo.On("Add", ptestArtwork).Return(nil)
// 		err = searcher.AddArtwork(ptestArtwork)

// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("failed artwork addition", func(t *testing.T) {
// 		// Arrange
// 		mockRepo := new(mockartworkrep.MockArtworkRep)
// 		searcher := &searcher{rep: mockRepo}
// 		var err error
// 		ptestArtwork, err := createTestArtwork()
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		expectedError := artworkrep.ErrFailedToAddArtwork

// 		mockRepo.On("Add", ptestArtwork).Return(expectedError)
// 		err = searcher.AddArtwork(ptestArtwork)

// 		assert.Error(t, err)
// 		assert.Equal(t, expectedError, err)
// 		mockRepo.AssertExpectations(t)
// 	})
// }
