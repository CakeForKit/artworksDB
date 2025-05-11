package artworkserv

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"github.com/google/uuid"
)

type ArtworkService interface {
	GetAllArtworks(ctx context.Context) ([]*models.Artwork, error)
	GetAllAuthors(ctx context.Context) ([]*models.Author, error)
	GetAllCollections(ctx context.Context) ([]*models.Collection, error)
	Add(ctx context.Context, aw *models.Artwork) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) error
}

var (
	ErrNoAuthor     = errors.New("failed to update the Artwork, no author")
	ErrNoCollection = errors.New("failed to update the Artwork, no collection")
)

type artworkService struct {
	artworkRep    artworkrep.ArtworkRep
	authorRep     authorrep.AuthorRep
	collectionRep collectionrep.CollectionRep
}

func NewArtworkService(artRep artworkrep.ArtworkRep, authorRep authorrep.AuthorRep, collectionRep collectionrep.CollectionRep) ArtworkService {
	return &artworkService{
		artworkRep:    artRep,
		authorRep:     authorRep,
		collectionRep: collectionRep,
	}
}

func (a *artworkService) GetAllArtworks(ctx context.Context) ([]*models.Artwork, error) {
	return a.artworkRep.GetAllArtworks(ctx)
}

func (a *artworkService) GetAllAuthors(ctx context.Context) ([]*models.Author, error) {
	return a.authorRep.GetAllAuthors(ctx)
}
func (a *artworkService) GetAllCollections(ctx context.Context) ([]*models.Collection, error) {
	return a.collectionRep.GetAllCollections(ctx)
}

func (a *artworkService) Add(ctx context.Context, aw *models.Artwork) error {
	authorExist, err := a.authorRep.CheckAuthorByID(ctx, aw.GetAuthor().GetID())
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w", err)
	}
	if !authorExist {
		err = a.authorRep.AddAuthor(ctx, aw.GetAuthor())
		if err != nil {
			return fmt.Errorf("artworkService.Add: %w", err)
		}
	}
	collectionExist, err := a.collectionRep.CheckCollectionByID(ctx, aw.GetCollection().GetID())
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w", err)
	}
	if !collectionExist {
		err = a.collectionRep.AddCollection(ctx, aw.GetCollection())
		if err != nil {
			return fmt.Errorf("artworkService.Add: %w", err)
		}
	}
	return a.artworkRep.Add(ctx, aw)
}

func (a *artworkService) Delete(ctx context.Context, id uuid.UUID) error {
	return a.artworkRep.Delete(ctx, id)
}

func (a *artworkService) Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Artwork) (*models.Artwork, error)) error {
	art, err := a.artworkRep.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("artworkService.Update: %w", err)
	}
	updatedArtwork, err := funcUpdate(art)
	if err != nil {
		return fmt.Errorf("artworkService.Update funcUpdate: %w", err)
	}
	updatedAuthorID := updatedArtwork.GetAuthor().GetID()
	updatedCollectionID := updatedArtwork.GetCollection().GetID()
	res, err := a.authorRep.CheckAuthorByID(ctx, updatedAuthorID)
	if err != nil {
		return fmt.Errorf("artworkService.Update: %w", err)
	} else if !res {
		err = a.authorRep.AddAuthor(ctx, updatedArtwork.GetAuthor())
		if err != nil {
			return fmt.Errorf("artworkService.Update: %w", err)
		}
	}
	res, err = a.collectionRep.CheckCollectionByID(ctx, updatedCollectionID)
	if err != nil {
		return fmt.Errorf("artworkService.Update: %w", err)
	} else if !res {
		err = a.collectionRep.AddCollection(ctx, updatedArtwork.GetCollection())
		if err != nil {
			return fmt.Errorf("artworkService.Update: %w", err)
		}
	}
	updateReq := artworkrep.ArtworkUpdate{
		Title:        updatedArtwork.GetTitle(),
		CreationYear: updatedArtwork.GetCreationYear(),
		Technic:      updatedArtwork.GetTechnic(),
		Material:     updatedArtwork.GetMaterial(),
		Size:         updatedArtwork.GetSize(),
		AuthorID:     updatedAuthorID,
		CollectionID: updatedCollectionID,
	}
	return a.artworkRep.Update(ctx, id, &updateReq)
}
