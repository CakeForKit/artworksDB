package artworkrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockArtworkRep struct {
	mock.Mock
}

func (m *MockArtworkRep) GetAllArtworks(ctx context.Context, filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error) {
	args := m.Called(ctx, filterOps, sortOps)
	return args.Get(0).([]*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Artwork, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) Add(ctx context.Context, aw *models.Artwork) error {
	args := m.Called(ctx, aw)
	return args.Error(0)
}

func (m *MockArtworkRep) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArtworkRep) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) error {
	args := m.Called(ctx, id, funcUpdate)
	return args.Error(0)
}

func (m *MockArtworkRep) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockArtworkRep) Close() {
	m.Called()
}
