package mockartworkrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
)

type MockArtworkRep struct {
	mock.Mock
}

func (m *MockArtworkRep) GetAll(ctx context.Context) ([]*models.Artwork, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Artwork, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) GetByTitle(ctx context.Context, title string) ([]*models.Artwork, error) {
	args := m.Called(ctx, title)
	return args.Get(0).([]*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) GetByAuthor(ctx context.Context, author *models.Author) ([]*models.Artwork, error) {
	args := m.Called(ctx, author)
	return args.Get(0).([]*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) GetByCreationTime(ctx context.Context, yearBeg int, yearEnd int) ([]*models.Artwork, error) {
	args := m.Called(ctx, yearBeg, yearEnd)
	return args.Get(0).([]*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) GetByEvent(ctx context.Context, event models.Event) ([]*models.Artwork, error) {
	args := m.Called(ctx, event)
	return args.Get(0).([]*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) Add(ctx context.Context, aw *models.Artwork) error {
	args := m.Called(ctx, aw)
	return args.Error(0)
}

func (m *MockArtworkRep) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArtworkRep) Update(ctx context.Context, id uuid.UUID,
	funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error) {

	args := m.Called(ctx, id, funcUpdate)
	return args.Get(0).(*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockArtworkRep) Close() {
}
