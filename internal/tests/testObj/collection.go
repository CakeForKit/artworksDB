package testobj

import (
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
)

type CollectionMother interface {
	CollectionP(collectionID uuid.UUID) *models.Collection
	CollectionUpdateReq() models.CollectionUpdateReq
}

func NewCollectionMother() CollectionMother {
	return &collectionMother{}
}

type collectionMother struct {
}

func (um *collectionMother) CollectionP(collectionID uuid.UUID) *models.Collection {
	collection, _ := models.NewCollection(
		collectionID,
		"test-collection-title",
	)
	return &collection
}

func (um *collectionMother) CollectionUpdateReq() models.CollectionUpdateReq {
	return models.CollectionUpdateReq{
		Title: "changed-test-collection-title",
	}
}
