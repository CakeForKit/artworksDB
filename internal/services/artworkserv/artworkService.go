package artworkserv

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"github.com/google/uuid"
)

type ArtworkService interface {
	GetAllArtworks(ctx context.Context) ([]*models.Artwork, error)
	Add(ctx context.Context, artworkReq jsonreqresp.AddArtworkRequest) error
	Delete(ctx context.Context, idArt uuid.UUID) error
	Update(ctx context.Context, idArt uuid.UUID, updateFields jsonreqresp.ArtworkUpdate) error
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
	return a.artworkRep.GetAllArtworks(ctx, &jsonreqresp.ArtworkFilter{}, &jsonreqresp.ArtworkSortOps{})
}

func (a *artworkService) Add(ctx context.Context, artworkReq jsonreqresp.AddArtworkRequest) error {
	author, err := a.authorRep.GetByID(ctx, uuid.MustParse(artworkReq.AuthorID))
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w", err)
	}
	collection, err := a.collectionRep.GetCollectionByID(ctx, uuid.MustParse(artworkReq.CollectionID))
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w", err)
	}

	artwork, err := models.NewArtwork(
		uuid.New(),
		artworkReq.Title,
		artworkReq.Technic,
		artworkReq.Material,
		artworkReq.Size,
		artworkReq.CreationYear,
		author,
		collection,
	)
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w: %w", models.ErrValidateArtwork, err)
	}

	err = a.artworkRep.Add(ctx, &artwork)
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w", err)
	}

	return nil
}

func (a *artworkService) Delete(ctx context.Context, idArt uuid.UUID) error {
	return a.artworkRep.Delete(ctx, idArt)
}

func (a *artworkService) Update(ctx context.Context, idArt uuid.UUID, updateFields jsonreqresp.ArtworkUpdate) error {
	_, err := a.authorRep.GetByID(ctx, uuid.MustParse(updateFields.AuthorID))
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w", err)
	}
	_, err = a.collectionRep.GetCollectionByID(ctx, uuid.MustParse(updateFields.CollectionID))
	if err != nil {
		return fmt.Errorf("artworkService.Add: %w", err)
	}
	return a.artworkRep.Update(
		ctx,
		idArt,
		func(a *models.Artwork) (*models.Artwork, error) {
			err := a.Update(updateFields)
			return a, err
		})
}
