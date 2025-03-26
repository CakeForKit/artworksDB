package artworkMemoryRep

import (
	"fmt"
	"sync"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/models"
	artworkRep "git.iu7.bmstu.ru/ped22u691/PPO.git/repository/artwork"
	"github.com/google/uuid"
)

type ArtworkMemoryRep struct {
	artworks map[uuid.UUID]models.Artwork
	sync.Mutex
}

func New() *ArtworkMemoryRep {
	return &ArtworkMemoryRep{artworks: make(map[uuid.UUID]models.Artwork)}
}

func (amr *ArtworkMemoryRep) Get(id uuid.UUID) (models.Artwork, error) {
	if artwork, ok := amr.artworks[id]; ok {
		return artwork, nil
	}
	return models.Artwork{}, artworkRep.ErrArtworkNotFound
}

func (amr *ArtworkMemoryRep) Add(aw models.Artwork) error {
	if amr.artworks == nil {
		amr.Lock()
		amr.artworks = make(map[uuid.UUID]models.Artwork)
		amr.Unlock()
	}
	// already in rep
	if _, ok := amr.artworks[aw.GetID()]; ok {
		return fmt.Errorf("artwork already exists: %w", artworkRep.ErrFailedToAddArtwork)
	}
	amr.Lock()
	amr.artworks[aw.GetID()] = aw
	amr.Unlock()
	return nil
}

func (amr *ArtworkMemoryRep) Update(aw models.Artwork) error {
	if _, ok := amr.artworks[aw.GetID()]; !ok {
		return fmt.Errorf("artwork doesn't exist: %w", artworkRep.ErrUpdateArtwork)
	}
	amr.Lock()
	amr.artworks[aw.GetID()] = aw
	amr.Unlock()
	return nil
}

func (amr *ArtworkMemoryRep) Delete(aw models.Artwork) error {
	var id uuid.UUID = aw.GetID()
	if _, ok := amr.artworks[id]; ok {
		amr.Lock()
		delete(amr.artworks, id)
		amr.Unlock()
	}
	return nil
}
