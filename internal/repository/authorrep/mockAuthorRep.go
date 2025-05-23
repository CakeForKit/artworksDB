package authorrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthorRep реализует AuthorRep интерфейс для тестирования
type MockAuthorRep struct {
	mock.Mock
}

func (m *MockAuthorRep) GetAll(ctx context.Context) ([]*models.Author, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Author), args.Error(1)
}

func (m *MockAuthorRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Author, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorRep) Add(ctx context.Context, a *models.Author) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAuthorRep) Delete(ctx context.Context, idAuthor uuid.UUID) error {
	args := m.Called(ctx, idAuthor)
	return args.Error(0)
}

func (m *MockAuthorRep) Update(ctx context.Context, idAuthor uuid.UUID, funcUpdate func(*models.Author) (*models.Author, error)) error {
	args := m.Called(ctx, idAuthor, funcUpdate)
	return args.Error(0)
}

func (m *MockAuthorRep) HasArtworks(ctx context.Context, authorID uuid.UUID) (bool, error) {
	args := m.Called(ctx, authorID)
	return args.Bool(0), args.Error(1)
}
