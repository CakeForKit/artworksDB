package searcher

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
)

type Searcher interface {
	GetAllArtworks(ctx context.Context, filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error)
	GetAllEvents(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error)
	// GetEventsOfArtworkOnDate(ctx context.Context, artworkID uuid.UUID, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error)
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

// func (s *searcher) GetEventsOfArtworkOnDate(ctx context.Context, artworkID uuid.UUID, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
// 	return s.eventRep.GetEventsOfArtworkOnDate(ctx, artworkID, dateBeg, dateEnd)
// }
