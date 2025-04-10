package artworkserv

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"github.com/google/uuid"
)

type ArtworkService interface {
	GetAllArtworks() []*models.Artwork
	Add(*models.Artwork) error
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error)
}

type artworkService struct {
	artworkRep artworkrep.ArtworkRep
}

func NewArtworkService(artRep artworkrep.ArtworkRep) ArtworkService {
	return &artworkService{
		artworkRep: artRep,
	}
}

func (a *artworkService) GetAllArtworks() []*models.Artwork {
	return a.artworkRep.GetAll()
}

func (a *artworkService) Add(aw *models.Artwork) error {
	return a.artworkRep.Add(aw)
}

func (a *artworkService) Delete(id uuid.UUID) error {
	return a.artworkRep.Delete(id)
}

func (a *artworkService) Update(id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error) {
	return a.artworkRep.Update(id, funcUpdate)
}
