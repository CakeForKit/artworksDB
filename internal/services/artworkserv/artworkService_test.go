package artworkserv_test

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep/mockartworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"github.com/google/uuid"
	"github.com/stateio/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestArtworkService_GetAllArtworks(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		expectedCount int
	}{
		{"single artwork", 1},
		{"multiple artworks", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockartworkrep.MockArtworkRep)
			service := artworkserv.NewArtworkService(mockRepo)

			var expectedArtworks []*models.Artwork
			for i := 0; i < tt.expectedCount; i++ {
				artwork := createTestArtwork()
				expectedArtworks = append(expectedArtworks, artwork)
			}

			mockRepo.On("GetAll", ctx).Return(expectedArtworks, nil)

			result, err := service.GetAllArtworks(ctx)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(result))
			assert.Equal(t, expectedArtworks, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestArtworkService_Add(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		mockError     error
		expectedError error
	}{
		{"success", nil, nil},
		{"repository error", artworkrep.ErrFailedToAddArtwork, artworkrep.ErrFailedToAddArtwork},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockartworkrep.MockArtworkRep)
			service := artworkserv.NewArtworkService(mockRepo)
			artwork := createTestArtwork()

			mockRepo.On("Add", ctx, artwork).Return(tt.mockError)

			err := service.Add(ctx, artwork)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestArtworkService_Delete(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		mockError     error
		expectedError error
	}{
		{"success", nil, nil},
		{"not found", artworkrep.ErrArtworkNotFound, artworkrep.ErrArtworkNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockartworkrep.MockArtworkRep)
			service := artworkserv.NewArtworkService(mockRepo)
			artworkID := uuid.New()

			mockRepo.On("Delete", ctx, artworkID).Return(tt.mockError)

			err := service.Delete(ctx, artworkID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestArtworkService_Update(t *testing.T) {
	ctx := context.Background()
	updateTitle := "Updated Title"
	updateFunc := func(aw *models.Artwork) (*models.Artwork, error) {
		newArtwork, err := models.NewArtwork(
			aw.GetID(),
			updateTitle,
			aw.GetTechnic(),
			aw.GetMaterial(),
			aw.GetSize(),
			aw.GetCreationYear(),
			aw.GetAuthor(),
			aw.GetCollection(),
		)
		if err != nil {
			return nil, err
		}
		return &newArtwork, nil
	}

	expectedArtwork, _ := updateFunc(createTestArtwork())
	tests := []struct {
		name          string
		mockResult    *models.Artwork
		mockError     error
		expectedError error
	}{
		{
			"success",
			expectedArtwork,
			nil,
			nil,
		},
		{
			"not found",
			nil,
			artworkrep.ErrArtworkNotFound,
			artworkrep.ErrArtworkNotFound,
		},
		{
			"update error",
			nil,
			artworkrep.ErrUpdateArtwork,
			artworkrep.ErrUpdateArtwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockartworkrep.MockArtworkRep)
			service := artworkserv.NewArtworkService(mockRepo)
			artworkID := uuid.New()

			mockRepo.On("Update", ctx, artworkID, mock.Anything).Return(tt.mockResult, tt.mockError)

			result, err := service.Update(ctx, artworkID, updateFunc)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, updateTitle, result.GetTitle())
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
