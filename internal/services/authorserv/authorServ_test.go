package authorserv_test

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/authorserv"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createTestAuthor() *models.Author {
	author, _ := models.NewAuthor(uuid.New(), "Test Author", 1900, 2000)
	return &author
}

func createTestUpdateRequest() models.AuthorUpdateReq {
	return models.AuthorUpdateReq{
		Name:      "Updated Author",
		BirthYear: 1901,
		DeathYear: 2001,
	}
}

func TestAuthorService_GetAll(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		setupMocks    func(*authorrep.MockAuthorRep)
		expectedCount int
		expectedError error
	}{
		{
			name: "success with authors",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				authors := []*models.Author{
					createTestAuthor(),
					createTestAuthor(),
				}
				m.On("GetAll", ctx).Return(authors, nil)
			},
			expectedCount: 2,
		},
		{
			name: "empty result",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("GetAll", ctx).Return([]*models.Author{}, nil)
			},
			expectedCount: 0,
		},
		{
			name: "repository error",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("GetAll", ctx).Return(nil, assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &authorrep.MockAuthorRep{}
			tt.setupMocks(mockRepo)

			service := authorserv.NewAuthorServ(mockRepo)
			result, err := service.GetAll(ctx)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthorService_Add(t *testing.T) {
	ctx := context.Background()
	testAuthor := createTestAuthor()

	tests := []struct {
		name          string
		setupMocks    func(*authorrep.MockAuthorRep)
		author        *models.Author
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("Add", ctx, testAuthor).Return(nil)
			},
			author: testAuthor,
		},
		{
			name: "repository error",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("Add", ctx, testAuthor).Return(assert.AnError)
			},
			author:        testAuthor,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &authorrep.MockAuthorRep{}
			tt.setupMocks(mockRepo)

			service := authorserv.NewAuthorServ(mockRepo)
			err := service.Add(ctx, tt.author)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthorService_Update(t *testing.T) {
	ctx := context.Background()
	authorID := uuid.New()
	testRequest := createTestUpdateRequest()

	tests := []struct {
		name          string
		setupMocks    func(*authorrep.MockAuthorRep)
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("Update", ctx, authorID, mock.Anything).Return(nil)
			},
		},
		{
			name: "repository error",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("Update", ctx, authorID, mock.Anything).Return(assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name: "validation error",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("Update", ctx, authorID, mock.Anything).
					Return(models.ErrAuthorBirthAfterDeath)
			},
			expectedError: models.ErrAuthorBirthAfterDeath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &authorrep.MockAuthorRep{}
			tt.setupMocks(mockRepo)

			service := authorserv.NewAuthorServ(mockRepo)
			err := service.Update(ctx, authorID, testRequest)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthorService_Delete(t *testing.T) {
	ctx := context.Background()
	authorID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func(*authorrep.MockAuthorRep)
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("HasArtworks", ctx, authorID).Return(false, nil)
				m.On("Delete", ctx, authorID).Return(nil)
			},
		},
		{
			name: "has linked artworks",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("HasArtworks", ctx, authorID).Return(true, nil)
			},
			expectedError: authorserv.ErrHasLinkedArtworks,
		},
		{
			name: "has artworks check error",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("HasArtworks", ctx, authorID).Return(false, assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name: "delete error",
			setupMocks: func(m *authorrep.MockAuthorRep) {
				m.On("HasArtworks", ctx, authorID).Return(false, nil)
				m.On("Delete", ctx, authorID).Return(assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &authorrep.MockAuthorRep{}
			tt.setupMocks(mockRepo)

			service := authorserv.NewAuthorServ(mockRepo)
			err := service.Delete(ctx, authorID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
