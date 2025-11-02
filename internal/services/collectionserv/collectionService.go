package collectionserv

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/google/uuid"
)

type CollectionServ interface {
	GetAll(ctx context.Context) ([]*models.Collection, error)
	Add(ctx context.Context, col *models.Collection) error
	Update(ctx context.Context, idCol uuid.UUID, updateReq models.CollectionUpdateReq) error
	Delete(ctx context.Context, idCol uuid.UUID) error
}

func NewCollectionServ(collectionRep collectionrep.CollectionRep) CollectionServ {
	return &collectionServ{
		collectionRep: collectionRep,
	}
}

type collectionServ struct {
	collectionRep collectionrep.CollectionRep
}

func (s *collectionServ) GetAll(ctx context.Context) ([]*models.Collection, error) {
	return s.collectionRep.GetAll(ctx)
}

func (s *collectionServ) Add(ctx context.Context, col *models.Collection) error {
	return s.collectionRep.Add(ctx, col)
}

func (s *collectionServ) Delete(ctx context.Context, idCol uuid.UUID) error {
	return s.collectionRep.Delete(ctx, idCol)
}

func (s *collectionServ) Update(
	ctx context.Context,
	idCol uuid.UUID,
	updateReq models.CollectionUpdateReq,
) error {
	return s.collectionRep.Update(
		ctx,
		idCol,
		func(c *models.Collection) (*models.Collection, error) {
			err := c.Update(updateReq)
			return c, err
		})
}
