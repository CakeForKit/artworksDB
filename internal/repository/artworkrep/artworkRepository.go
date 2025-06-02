package artworkrep

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

var (
	ErrArtworkNotFound = errors.New("the Artwork was not found in the repository")
	ErrUpdateArtwork   = errors.New("err update artwork params")
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

func NewArtworkRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (ArtworkRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgArtworkRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHArtworkRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewArtworkRep: %w", cnfg.ErrUnknownDB)
	}
	// return &MockAdminRep{}, nil
}
