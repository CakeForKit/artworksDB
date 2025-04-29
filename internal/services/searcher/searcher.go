package searcher

import (
	"context"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
)

type Searcher interface {
	GetAllArtworks(ctx context.Context) ([]*models.Artwork, error)
	// GetByID(id uuid.UUID) (*models.Artwork, error)
	// AddArtwork(*models.Artwork) error
	FilterArtworkByTitle(ctx context.Context, title string) ([]*models.Artwork, error)
	FilterArtworkByAuthor(ctx context.Context, author *models.Author) ([]*models.Artwork, error)
	FilterArtworkByCreationTime(ctx context.Context, yearBeg int, yearEnd int) ([]*models.Artwork, error)
	FilterArtworkByEvent(ctx context.Context, event models.Event) ([]*models.Artwork, error)
	GetAllEvents(ctx context.Context) ([]*models.Event, error)
	FilterEventsByDate(ctx context.Context, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error)
	GetEventOfArtworkOnDate(ctx context.Context, artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error)
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

func (s *searcher) GetAllArtworks(ctx context.Context) ([]*models.Artwork, error) {
	return s.artworkRep.GetAll(ctx)
}

// func (s *searcher) GetByID(id uuid.UUID) (*models.Artwork, error) {
// 	return s.artowrkRep.GetByID(id)
// }

func (s *searcher) FilterArtworkByTitle(ctx context.Context, title string) ([]*models.Artwork, error) {
	return s.artworkRep.GetByTitle(ctx, title)
}

func (s *searcher) FilterArtworkByAuthor(ctx context.Context, author *models.Author) ([]*models.Artwork, error) {
	return s.artworkRep.GetByAuthor(ctx, author)
}

func (s *searcher) FilterArtworkByCreationTime(ctx context.Context, yearBeg int, yearEnd int) ([]*models.Artwork, error) {
	return s.artworkRep.GetByCreationTime(ctx, yearBeg, yearEnd)
}

func (s *searcher) FilterArtworkByEvent(ctx context.Context, event models.Event) ([]*models.Artwork, error) {
	return s.artworkRep.GetByEvent(ctx, event)
}

func (s *searcher) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	return s.eventRep.GetAll(ctx)
}

func (s *searcher) FilterEventsByDate(ctx context.Context, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	return s.eventRep.GetByDate(ctx, dateBeg, dateEnd)
}

func (s *searcher) GetEventOfArtworkOnDate(ctx context.Context, artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error) {
	return s.eventRep.GetEventOfArtworkOnDate(ctx, artwork, dateBeg, dateEnd)
}
