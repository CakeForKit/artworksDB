package authorrep

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrAuthorNotFound = errors.New("the Author was not found in the repository")
)

type AuthorRep interface {
	GetAllAuthors(ctx context.Context) ([]*models.Author, error)
	CheckAuthorByID(ctx context.Context, id uuid.UUID) (bool, error)
	AddAuthor(ctx context.Context, e *models.Author) error
}

func NewAuthorRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (AuthorRep, error) {
	rep, err := NewPgAuthorRep(ctx, pgCreds, dbConf)
	return rep, err
}
