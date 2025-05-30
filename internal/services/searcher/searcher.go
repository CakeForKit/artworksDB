package searcher

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"github.com/google/uuid"
)

type Searcher interface {
	GetAllArtworks(ctx context.Context, filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error)
	GetAllEvents(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error)
	GetEvent(ctx context.Context, eventID uuid.UUID) (*models.Event, error)
	GetArtworksFromEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Artwork, error)
	GetCollectionsStat(ctx context.Context, eventID uuid.UUID) ([]*models.StatCollections, error)
}

type searcher struct {
	artworkRep artworkrep.ArtworkRep
	eventRep   eventrep.EventRep
}

func NewSearcher(artRep artworkrep.ArtworkRep, eventRep eventrep.EventRep) Searcher {
	return &searcher{
		artworkRep: artRep,
		eventRep:   eventRep,
	}
}

func (s *searcher) GetAllArtworks(ctx context.Context, filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error) {
	return s.artworkRep.GetAllArtworks(ctx, filterOps, sortOps)
}

func (s *searcher) GetAllEvents(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error) {
	if filterOps.DateBegin.After(filterOps.DateEnd) {
		return nil, fmt.Errorf("searcher.GetAllEvents : %w", jsonreqresp.ErrEventFilterDate)
	}
	return s.eventRep.GetAll(ctx, filterOps)
}

func (s *searcher) GetEvent(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	return s.eventRep.GetByID(ctx, eventID)
}

func (s *searcher) GetArtworksFromEvent(ctx context.Context, eventID uuid.UUID) ([]*models.Artwork, error) {
	artworkIDs, err := s.eventRep.GetArtworkIDs(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("searcher.GetArtworkFromEvent: %w", err)
	}
	artworks := make([]*models.Artwork, len(artworkIDs))
	for i, aID := range artworkIDs {
		art, err := s.artworkRep.GetByID(ctx, aID)
		if err != nil {
			return nil, fmt.Errorf("searcher.GetArtworkFromEvent: %w", err)
		}
		artworks[i] = art
	}
	return artworks, nil
}

func (s *searcher) GetCollectionsStat(ctx context.Context, eventID uuid.UUID) ([]*models.StatCollections, error) {
	_, err := s.eventRep.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("searcher.GetCollectionsStat: %w", err)
	}
	statCols, err := s.eventRep.GetCollectionsStat(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("searcher.GetCollectionsStat: %w", err)
	}
	return statCols, nil
}
