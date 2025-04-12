package searcher

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
)

type Searcher interface {
	GetAllArtworks() []*models.Artwork
	// GetByID(id uuid.UUID) (*models.Artwork, error)
	// AddArtwork(*models.Artwork) error
	FilterArtworkByTitle(title string) []*models.Artwork
	FilterArtworkByAuthor(author *models.Author) []*models.Artwork
	FilterArtworkByCreationTime(yearBeg int, yearEnd int) []*models.Artwork
	FilterArtworkByEvent(event models.Event) []*models.Artwork
	GetAllEvents() []*models.Event
	FilterEventsByDate(dateBeg time.Time, dateEnd time.Time) []*models.Event
	GetEventOfArtworkOnDate(artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error)
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

func (s *searcher) GetAllArtworks() []*models.Artwork {
	return s.artworkRep.GetAll()
}

// func (s *searcher) GetByID(id uuid.UUID) (*models.Artwork, error) {
// 	return s.artowrkRep.GetByID(id)
// }

func (s *searcher) FilterArtworkByTitle(title string) []*models.Artwork {
	return s.artworkRep.GetByTitle(title)
}

func (s *searcher) FilterArtworkByAuthor(author *models.Author) []*models.Artwork {
	return s.artworkRep.GetByAuthor(author)
}

func (s *searcher) FilterArtworkByCreationTime(yearBeg int, yearEnd int) []*models.Artwork {
	return s.artworkRep.GetByCreationTime(yearBeg, yearEnd)
}

func (s *searcher) FilterArtworkByEvent(event models.Event) []*models.Artwork {
	return s.artworkRep.GetByEvent(event)
}

func (s *searcher) GetAllEvents() []*models.Event {
	return s.eventRep.GetAll()
}

func (s *searcher) FilterEventsByDate(dateBeg time.Time, dateEnd time.Time) []*models.Event {
	return s.eventRep.GetByDate(dateBeg, dateEnd)
}

func (s *searcher) GetEventOfArtworkOnDate(artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) (*models.Event, error) {
	return s.eventRep.GetEventOfArtworkOnDate(artwork, dateBeg, dateEnd)
}
