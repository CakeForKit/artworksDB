package artworkrep

import (
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep/mockartworkrep"
	"github.com/google/uuid"
)

var (
	ErrArtworkNotFound    = errors.New("the Artwork was not found in the repository")
	ErrFailedToAddArtwork = errors.New("failed to add the Artwork to the repository")
	ErrUpdateArtwork      = errors.New("failed to update the Artwork in the repository")
)

type ArtworkRep interface {
	GetAll() []*models.Artwork
	GetByID(uuid.UUID) (*models.Artwork, error)
	GetByTitle(title string) []*models.Artwork
	GetByAuthor(author *models.Author) []*models.Artwork
	GetByCreationTime(yearBeg int, yearEnd int) []*models.Artwork
	GetByEvent(event models.Event) []*models.Artwork
	//
	Add(aw *models.Artwork) error
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error)
}

func NewArtworkRep() ArtworkRep {
	return &mockartworkrep.MockArtworkRep{}
}
