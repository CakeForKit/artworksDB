package testobj

import (
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type ArtworkMother interface {
	ArtworkP(artworkID uuid.UUID) *models.Artwork
	AddArtworkRequest(creationYear int, authorID uuid.UUID, collectionID uuid.UUID) jsonreqresp.AddArtworkRequest
	ArtworkFilter() jsonreqresp.ArtworkFilter
	ArtworkSortOps() jsonreqresp.ArtworkSortOps
}

func NewArtworkMother() ArtworkMother {
	return &artworkMother{}
}

type artworkMother struct {
}

func (um *artworkMother) ArtworkP(artworkID uuid.UUID) *models.Artwork {
	authorCretor := NewAuthorMother()
	colsCreator := NewCollectionMother()
	artwork, _ := models.NewArtwork(
		artworkID,
		"test-artwork-title",
		"test-artwork-technic",
		"test-artwork-material",
		"test-artwork-size",
		1999,
		authorCretor.AuthorBirthYearP(uuid.New(), 1999-5),
		colsCreator.CollectionP(uuid.New()),
	)
	return &artwork
}

func (um *artworkMother) AddArtworkRequest(creationYear int, authorID uuid.UUID, collectionID uuid.UUID) jsonreqresp.AddArtworkRequest {
	return jsonreqresp.AddArtworkRequest{
		Title:        "test-artwork-title",
		CreationYear: creationYear,
		Technic:      "test-artwork-technic",
		Material:     "test-artwork-material",
		Size:         "test-artwork-size",
		AuthorID:     authorID.String(),
		CollectionID: collectionID.String(),
	}
}

func (um *artworkMother) ArtworkFilter() jsonreqresp.ArtworkFilter {
	return jsonreqresp.ArtworkFilter{
		Title:      "test-artwork-title",
		AuthorName: "test-artwork-author",
		Collection: "test-artwork-collection",
		EventID:    uuid.New(),
	}
}

func (um *artworkMother) ArtworkSortOps() jsonreqresp.ArtworkSortOps {
	return jsonreqresp.ArtworkSortOps{
		Field:     jsonreqresp.TitleSortFieldArtwork,
		Direction: jsonreqresp.ASCDirection,
	}
}
