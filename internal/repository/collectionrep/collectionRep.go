package collectionrep

import (
	"context"
	"errors"
	"fmt"

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

func NewCollectionRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (CollectionRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgCollectionRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHCollectionRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewCollectionRep: %w", cnfg.ErrUnknownDB)
	}
	// return &MockAdminRep{}, nil
}
