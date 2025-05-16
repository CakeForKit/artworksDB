package artworkrep

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

var (
	ErrArtworkNotFound    = errors.New("the Artwork was not found in the repository")
	ErrFailedToAddArtwork = errors.New("failed to add the Artwork to the repository")
	ErrUpdateArtwork      = errors.New("err update artwork params")
)

type ArtworkRep interface {
	GetAllArtworks(ctx context.Context, filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Artwork, error)
	//
	Add(ctx context.Context, aw *models.Artwork) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) error
	Ping(ctx context.Context) error
	Close()
}

func NewArtworkRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (ArtworkRep, error) {
	rep, err := NewPgArtworkRep(ctx, pgCreds, dbConf)
	return rep, err
}
