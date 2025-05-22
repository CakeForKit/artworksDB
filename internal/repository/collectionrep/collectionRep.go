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
	ErrUpdateCollection   = errors.New("err update Collection params")
)

type CollectionRep interface {
	GetAllCollections(ctx context.Context) ([]*models.Collection, error)
	GetCollectionByID(ctx context.Context, id uuid.UUID) (*models.Collection, error)
	AddCollection(ctx context.Context, e *models.Collection) error
	DeleteCollection(ctx context.Context, idCol uuid.UUID) error
	UpdateCollection(ctx context.Context, idCol uuid.UUID, funcUpdate func(*models.Collection) (*models.Collection, error)) error
}

func NewCollectionRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (CollectionRep, error) {
	return NewPgCollectionRep(ctx, pgCreds, dbConf)
	// return &MockCollectionRep{}, nil
}
