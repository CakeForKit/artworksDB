package artworkserv_test

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func createTestAddRequest(authorID, collectionID uuid.UUID) jsonreqresp.AddArtworkRequest {
	return jsonreqresp.AddArtworkRequest{
		Title:        "Test Artwork",
		Technic:      "oil on canvas",
		Material:     "canvas",
		Size:         "100x100 cm",
		CreationYear: 1950,
		AuthorID:     authorID.String(),
		CollectionID: collectionID.String(),
	}
}

func createTestUpdateRequest(authorID, collectionID uuid.UUID) jsonreqresp.ArtworkUpdate {
	return jsonreqresp.ArtworkUpdate{
		Title:        "Updated Artwork",
		Technic:      "updated technic",
		Material:     "updated material",
		Size:         "200x200 cm",
		CreationYear: 2000,
		AuthorID:     authorID.String(),
		CollectionID: collectionID.String(),
	}
}

func TestArtworkService_GetAllArtworks(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		setupMocks    func(*artworkrep.MockArtworkRep)
		expectedCount int
		expectedError error
	}{
		{
			name: "success with artworks",
			setupMocks: func(m *artworkrep.MockArtworkRep) {
				artworks := []*models.Artwork{
					createTestArtwork(),
					createTestArtwork(),
				}
				m.On("GetAllArtworks", ctx, &jsonreqresp.ArtworkFilter{}, &jsonreqresp.ArtworkSortOps{}).
					Return(artworks, nil)
			},
			expectedCount: 2,
		},
		{
			name: "empty result",
			setupMocks: func(m *artworkrep.MockArtworkRep) {
				m.On("GetAllArtworks", ctx, &jsonreqresp.ArtworkFilter{}, &jsonreqresp.ArtworkSortOps{}).
					Return([]*models.Artwork{}, nil)
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artMock := &artworkrep.MockArtworkRep{}
			authMock := &authorrep.MockAuthorRep{}
			colMock := &collectionrep.MockCollectionRep{}

			tt.setupMocks(artMock)

			service := artworkserv.NewArtworkService(artMock, authMock, colMock)
			result, err := service.GetAllArtworks(ctx)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result))
			}

			artMock.AssertExpectations(t)
		})
	}
}

func TestArtworkService_Add(t *testing.T) {
	ctx := context.Background()
	testAuthor := createTestAuthor()
	testCollection := createTestCollection()
	testRequest := createTestAddRequest(testAuthor.GetID(), testCollection.GetID())

	tests := []struct {
		name          string
		setupMocks    func(*artworkrep.MockArtworkRep, *authorrep.MockAuthorRep, *collectionrep.MockCollectionRep)
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func(art *artworkrep.MockArtworkRep, auth *authorrep.MockAuthorRep, col *collectionrep.MockCollectionRep) {
				auth.On("GetByID", ctx, testAuthor.GetID()).Return(testAuthor, nil)
				col.On("GetCollectionByID", ctx, testCollection.GetID()).Return(testCollection, nil)
				art.On("Add", ctx, mock.AnythingOfType("*models.Artwork")).Return(nil)
			},
		},
		{
			name: "author not found",
			setupMocks: func(art *artworkrep.MockArtworkRep, auth *authorrep.MockAuthorRep, col *collectionrep.MockCollectionRep) {
				auth.On("GetByID", ctx, testAuthor.GetID()).Return(nil, authorrep.ErrAuthorNotFound)
			},
			expectedError: authorrep.ErrAuthorNotFound,
		},
		{
			name: "collection not found",
			setupMocks: func(art *artworkrep.MockArtworkRep, auth *authorrep.MockAuthorRep, col *collectionrep.MockCollectionRep) {
				auth.On("GetByID", ctx, testAuthor.GetID()).Return(testAuthor, nil)
				col.On("GetCollectionByID", ctx, testCollection.GetID()).Return(nil, collectionrep.ErrCollectionNotFound)
			},
			expectedError: collectionrep.ErrCollectionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artMock := &artworkrep.MockArtworkRep{}
			authMock := &authorrep.MockAuthorRep{}
			colMock := &collectionrep.MockCollectionRep{}

			tt.setupMocks(artMock, authMock, colMock)

			service := artworkserv.NewArtworkService(artMock, authMock, colMock)
			err := service.Add(ctx, testRequest)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			artMock.AssertExpectations(t)
			authMock.AssertExpectations(t)
			colMock.AssertExpectations(t)
		})
	}
}

func TestArtworkService_Delete(t *testing.T) {
	ctx := context.Background()
	artworkID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func(*artworkrep.MockArtworkRep)
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func(m *artworkrep.MockArtworkRep) {
				m.On("Delete", ctx, artworkID).Return(nil)
			},
		},
		{
			name: "not found",
			setupMocks: func(m *artworkrep.MockArtworkRep) {
				m.On("Delete", ctx, artworkID).Return(artworkrep.ErrArtworkNotFound)
			},
			expectedError: artworkrep.ErrArtworkNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artMock := &artworkrep.MockArtworkRep{}
			authMock := &authorrep.MockAuthorRep{}
			colMock := &collectionrep.MockCollectionRep{}

			tt.setupMocks(artMock)

			service := artworkserv.NewArtworkService(artMock, authMock, colMock)
			err := service.Delete(ctx, artworkID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			artMock.AssertExpectations(t)
		})
	}
}

func TestArtworkService_Update(t *testing.T) {
	ctx := context.Background()
	artworkID := uuid.New()
	testAuthor := createTestAuthor()
	testCollection := createTestCollection()
	testRequest := createTestUpdateRequest(testAuthor.GetID(), testCollection.GetID())

	tests := []struct {
		name          string
		setupMocks    func(*artworkrep.MockArtworkRep, *authorrep.MockAuthorRep, *collectionrep.MockCollectionRep)
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func(art *artworkrep.MockArtworkRep, auth *authorrep.MockAuthorRep, col *collectionrep.MockCollectionRep) {
				auth.On("GetByID", ctx, testAuthor.GetID()).Return(testAuthor, nil)
				col.On("GetCollectionByID", ctx, testCollection.GetID()).Return(testCollection, nil)
				art.On("Update", ctx, artworkID, mock.Anything).Return(nil)
			},
		},
		{
			name: "author not found",
			setupMocks: func(art *artworkrep.MockArtworkRep, auth *authorrep.MockAuthorRep, col *collectionrep.MockCollectionRep) {
				auth.On("GetByID", ctx, testAuthor.GetID()).Return(nil, authorrep.ErrAuthorNotFound)
			},
			expectedError: authorrep.ErrAuthorNotFound,
		},
		{
			name: "collection not found",
			setupMocks: func(art *artworkrep.MockArtworkRep, auth *authorrep.MockAuthorRep, col *collectionrep.MockCollectionRep) {
				auth.On("GetByID", ctx, testAuthor.GetID()).Return(testAuthor, nil)
				col.On("GetCollectionByID", ctx, testCollection.GetID()).Return(nil, collectionrep.ErrCollectionNotFound)
			},
			expectedError: collectionrep.ErrCollectionNotFound,
		},
		{
			name: "update error",
			setupMocks: func(art *artworkrep.MockArtworkRep, auth *authorrep.MockAuthorRep, col *collectionrep.MockCollectionRep) {
				auth.On("GetByID", ctx, testAuthor.GetID()).Return(testAuthor, nil)
				col.On("GetCollectionByID", ctx, testCollection.GetID()).Return(testCollection, nil)
				art.On("Update", ctx, artworkID, mock.Anything).Return(artworkrep.ErrUpdateArtwork)
			},
			expectedError: artworkrep.ErrUpdateArtwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artMock := &artworkrep.MockArtworkRep{}
			authMock := &authorrep.MockAuthorRep{}
			colMock := &collectionrep.MockCollectionRep{}

			tt.setupMocks(artMock, authMock, colMock)

			service := artworkserv.NewArtworkService(artMock, authMock, colMock)
			err := service.Update(ctx, artworkID, testRequest)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			artMock.AssertExpectations(t)
			authMock.AssertExpectations(t)
			colMock.AssertExpectations(t)
		})
	}
}
