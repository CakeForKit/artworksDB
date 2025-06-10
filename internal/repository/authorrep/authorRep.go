package authorrep

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrAuthorNotFound = errors.New("the Author was not found in the repository")
	ErrUpdateAuthor   = errors.New("err update author params")
)

type AuthorRep interface {
	GetAll(ctx context.Context) ([]*models.Author, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Author, error)
	Add(ctx context.Context, a *models.Author) error
	Delete(ctx context.Context, idAuthor uuid.UUID) error
	Update(ctx context.Context, idAuthor uuid.UUID, funcUpdate func(*models.Author) (*models.Author, error)) error
	HasArtworks(ctx context.Context, authorID uuid.UUID) (bool, error)
}

func NewAuthorRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (AuthorRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgAuthorRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHAuthorRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewAuthorRep: %w", cnfg.ErrUnknownDB)
	}
	// return &MockAdminRep{}, nil
}
