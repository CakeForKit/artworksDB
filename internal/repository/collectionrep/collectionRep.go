package collectionrep

import (
	"context"
	"errors"
	"fmt"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
)

var (
	ErrCollectionNotFound = errors.New("the Collection was not found in the repository")
	ErrUpdate             = errors.New("err update Collection params")
)

type CollectionRep interface {
	GetAll(ctx context.Context) ([]*models.Collection, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Collection, error)
	Add(ctx context.Context, e *models.Collection) error
	Delete(ctx context.Context, idCol uuid.UUID) error
	Update(ctx context.Context, idCol uuid.UUID, funcUpdate func(*models.Collection) (*models.Collection, error)) error
	Ping(ctx context.Context) error
	Close()
}

func NewCollectionRep(ctx context.Context, datebaseType string, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (CollectionRep, error) {
	if datebaseType == cnfg.PostgresDB {
		return NewPgCollectionRep(ctx, pgCreds, dbConf)
	} else if datebaseType == cnfg.ClickHouseDB {
		return NewCHCollectionRep(ctx, (*cnfg.ClickHouseCredentials)(pgCreds), dbConf)
	} else {
		return nil, fmt.Errorf("NewCollectionRep: %w", cnfg.ErrUnknownDB)
	}
}
