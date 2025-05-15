package artworkrep

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrArtworkNotFound    = errors.New("the Artwork was not found in the repository")
	ErrFailedToAddArtwork = errors.New("failed to add the Artwork to the repository")
	ErrUpdateArtwork      = errors.New("failed to update the Artwork in the repository")
	ErrNoAuthor           = errors.New("failed to update the Artwork, no author")
	ErrNoCollection       = errors.New("failed to update the Artwork, no collection")
)

type ArtworkRep interface {
	GetAll(ctx context.Context) ([]*models.Artwork, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Artwork, error)
	GetByTitle(ctx context.Context, title string) ([]*models.Artwork, error)
	GetByAuthor(ctx context.Context, author *models.Author) ([]*models.Artwork, error)
	GetByCreationTime(ctx context.Context, yearBeg int, yearEnd int) ([]*models.Artwork, error)
	GetByEvent(ctx context.Context, event models.Event) ([]*models.Artwork, error)
	//
	Add(ctx context.Context, aw *models.Artwork) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error)
	Ping(ctx context.Context) error
	Close()
}

func NewArtworkRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (ArtworkRep, error) {
	rep, err := NewPgArtworkRep(ctx, pgCreds, dbConf)
	return rep, err
}
