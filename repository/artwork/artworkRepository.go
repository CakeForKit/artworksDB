package artworkRep

import (
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/models"
	"github.com/google/uuid"
)

var (
	ErrArtworkNotFound    = errors.New("the Artwork was not found in the repository")
	ErrFailedToAddArtwork = errors.New("failed to add the Artwork to the repository")
	ErrUpdateArtwork      = errors.New("failed to update the Artwork in the repository")
)

type ArtworkRepository interface {
	Get(uuid.UUID) (models.Artwork, error)
	Add(models.Artwork) error
	Update(models.Artwork) error
	Delete(models.Artwork) error
}
