package mockartworkrep

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
	"github.com/stateio/testify/mock"
)

type MockArtworkRep struct {
	mock.Mock
}

func (m *MockArtworkRep) GetAll() []*models.Artwork {
	args := m.Called()
	return args.Get(0).([]*models.Artwork)
}

func (m *MockArtworkRep) GetByID(id uuid.UUID) (*models.Artwork, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Artwork), args.Error(1)
}

func (m *MockArtworkRep) GetByTitle(title string) []*models.Artwork {
	args := m.Called(title)
	return args.Get(0).([]*models.Artwork)
}

func (m *MockArtworkRep) GetByAuthor(author *models.Author) []*models.Artwork {
	args := m.Called(author)
	return args.Get(0).([]*models.Artwork)
}

func (m *MockArtworkRep) GetByCreationTime(yearBeg int, yearEnd int) []*models.Artwork {
	args := m.Called(yearBeg, yearEnd)
	return args.Get(0).([]*models.Artwork)
}

func (m *MockArtworkRep) GetByEvent(event models.Event) []*models.Artwork {
	args := m.Called(event)
	return args.Get(0).([]*models.Artwork)
}

func (m *MockArtworkRep) Add(aw *models.Artwork) error {
	args := m.Called(aw)
	return args.Error(0)
}

func (m *MockArtworkRep) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockArtworkRep) Update(id uuid.UUID,
	funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error) {

	args := m.Called(id, funcUpdate)
	return args.Get(0).(*models.Artwork), args.Error(1)
}
