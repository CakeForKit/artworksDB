package artworkserv

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"github.com/google/uuid"
)

type ArtworkService interface {
	GetAllArtworks(ctx context.Context) ([]*models.Artwork, error)
	Add(ctx context.Context, aw *models.Artwork) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error)
}

type artworkService struct {
	artworkRep artworkrep.ArtworkRep
}

func NewArtworkService(artRep artworkrep.ArtworkRep) ArtworkService {
	return &artworkService{
		artworkRep: artRep,
	}
}

func (a *artworkService) GetAllArtworks(ctx context.Context) ([]*models.Artwork, error) {
	return a.artworkRep.GetAll(ctx)
}

func (a *artworkService) Add(ctx context.Context, aw *models.Artwork) error {
	return a.artworkRep.Add(ctx, aw)
}

func (a *artworkService) Delete(ctx context.Context, id uuid.UUID) error {
	return a.artworkRep.Delete(ctx, id)
}

func (a *artworkService) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error) {
	return a.artworkRep.Update(ctx, id, funcUpdate)
}
