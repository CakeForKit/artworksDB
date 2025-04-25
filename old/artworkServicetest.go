package artworkserv_test

import (
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep/mockartworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"github.com/google/uuid"
	"github.com/stateio/testify/assert"
	"github.com/stateio/testify/mock"
)

func createTestAuthor() models.Author {
	author, _ := models.NewAuthor(uuid.New(), "Test Author", 1900, 2000)
	return author
}

func createTestCollection() models.Collection {
	collection, _ := models.NewCollection(uuid.New(), "Test Collection")
	return collection
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
		&author,
		&collection,
	)
	return &artwork
}

func TestArtworkService_GetAllArtworks(t *testing.T) {
	mockArt := new(mockartworkrep.MockArtworkRep)
	artworkService := artworkserv.NewArtworkService(mockArt)

	artwork := createTestArtwork()
	expectedArtworks := []*models.Artwork{&artwork}

	mockArt.On("GetAll").Return(expectedArtworks)

	t.Run("get all artworks", func(t *testing.T) {
		result := artworkService.GetAllArtworks()
		assert.Equal(t, expectedArtworks, result)
		assert.Len(t, result, 1)
		assert.Equal(t, artwork.GetTitle(), result[0].GetTitle())
		mockArt.AssertExpectations(t)
	})
}

func TestArtworkService_Add(t *testing.T) {

	artwork := createTestArtwork()

	t.Run("successful artwork addition", func(t *testing.T) {
		mockArt := new(mockartworkrep.MockArtworkRep)
		artworkService := artworkserv.NewArtworkService(mockArt)
		mockArt.On("Add", &artwork).Return(nil)

		err := artworkService.Add(&artwork)
		assert.NoError(t, err)
		mockArt.AssertExpectations(t)
	})

	t.Run("failed artwork addition", func(t *testing.T) {
		mockArt := new(mockartworkrep.MockArtworkRep)
		artworkService := artworkserv.NewArtworkService(mockArt)
		expectedError := artworkrep.ErrFailedToAddArtwork
		mockArt.On("Add", &artwork).Return(expectedError)

		err := artworkService.Add(&artwork)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockArt.AssertExpectations(t)
	})
}

func TestArtworkService_Delete(t *testing.T) {
	artwork := createTestArtwork()
	artworkID := artwork.GetID()

	t.Run("successful artwork deletion", func(t *testing.T) {
		mockArt := new(mockartworkrep.MockArtworkRep)
		artworkService := artworkserv.NewArtworkService(mockArt)
		mockArt.On("Delete", artworkID).Return(nil)

		err := artworkService.Delete(artworkID)
		assert.NoError(t, err)
		mockArt.AssertExpectations(t)
	})

	t.Run("failed artwork deletion - not found", func(t *testing.T) {
		mockArt := new(mockartworkrep.MockArtworkRep)
		artworkService := artworkserv.NewArtworkService(mockArt)
		expectedError := artworkrep.ErrArtworkNotFound
		mockArt.On("Delete", artworkID).Return(expectedError)

		err := artworkService.Delete(artworkID)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockArt.AssertExpectations(t)
	})
}

func TestArtworkService_Update(t *testing.T) {

	artwork := createTestArtwork()
	artworkID := artwork.GetID()
	newTitle := "Updated Title"
	var updateFunc func(aw *models.Artwork) (*models.Artwork, error) = func(aw *models.Artwork) (*models.Artwork, error) {
		updatedArtwork, err := models.NewArtwork(
			aw.GetID(),
			newTitle,
			aw.GetTechnic(),
			aw.GetMaterial(),
			aw.GetSize(),
			aw.GetCreationYear(),
			aw.GetAuthor(),
			aw.GetCollection(),
		)
		assert.NoError(t, err)
		return &updatedArtwork, nil
	}

	t.Run("successful artwork update", func(t *testing.T) {
		mockArt := new(mockartworkrep.MockArtworkRep)
		artworkService := artworkserv.NewArtworkService(mockArt)
		mockArt.On("Update", artworkID, mock.AnythingOfType("func(*models.Artwork) (*models.Artwork, error)")).Return(&artwork, nil)

		updatedArtwork, err := artworkService.Update(artworkID, updateFunc)
		assert.NoError(t, err)
		assert.Equal(t, artwork.GetID(), updatedArtwork.GetID())
		mockArt.AssertExpectations(t)
	})

	t.Run("failed artwork update - not found", func(t *testing.T) {
		mockArt := new(mockartworkrep.MockArtworkRep)
		artworkService := artworkserv.NewArtworkService(mockArt)
		expectedError := artworkrep.ErrArtworkNotFound
		var awnil *models.Artwork = nil
		mockArt.On("Update", artworkID, mock.AnythingOfType("func(*models.Artwork) (*models.Artwork, error)")).Return(awnil, artworkrep.ErrArtworkNotFound)

		updatedArtwork, err := artworkService.Update(artworkID, updateFunc)
		assert.Error(t, err)
		assert.Nil(t, updatedArtwork)
		assert.Equal(t, expectedError, err)
		mockArt.AssertExpectations(t)
	})

	t.Run("failed artwork update - update error", func(t *testing.T) {
		mockArt := new(mockartworkrep.MockArtworkRep)
		artworkService := artworkserv.NewArtworkService(mockArt)
		expectedError := artworkrep.ErrUpdateArtwork
		var awnil *models.Artwork = nil
		mockArt.On("Update", artworkID, mock.AnythingOfType("func(*models.Artwork) (*models.Artwork, error)")).Return(awnil, expectedError)

		updatedArtwork, err := artworkService.Update(artworkID, updateFunc)
		assert.Error(t, err)
		assert.Nil(t, updatedArtwork)
		assert.Equal(t, expectedError, err)
		mockArt.AssertExpectations(t)
	})
}
