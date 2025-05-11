package collectionrep

import (
	"context"
	"errors"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrCollectionNotFound = errors.New("the Collection was not found in the repository")
)

type CollectionRep interface {
	GetAllCollections(ctx context.Context) ([]*models.Collection, error)
	CheckCollectionByID(ctx context.Context, id uuid.UUID) (bool, error)
	AddCollection(ctx context.Context, e *models.Collection) error
}

func NewCollectionRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (CollectionRep, error) {
	rep, err := NewPgCollectionRep(ctx, pgCreds, dbConf)
	return rep, err
}
