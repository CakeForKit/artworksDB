package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Artwork struct {
	id           uuid.UUID
	title        string
	creationYear int
	technic      string
	material     string
	size         string
	author       *Author
	collection   *Collection
}

type ArtworkResponse struct {
	ID           string             `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title        string             `json:"title" example:"Mona Lisa"`
	CreationYear int                `json:"creationYear" example:"1503"`
	Technic      string             `json:"technic" example:"Oil painting"`
	Material     string             `json:"material" example:"Poplar wood"`
	Size         string             `json:"size" example:"77 cm × 53 cm"`
	Author       AuthorResponse     `json:"author"`
	Collection   CollectionResponse `json:"collection"`
}

type CreateArtworkRequest struct {
	Title        string `json:"title" binding:"required"`
	CreationYear int    `json:"creationYear" binding:"required"`
	Technic      string `json:"technic" binding:"required"`
	Material     string `json:"material" binding:"required"`
	Size         string `json:"size" binding:"required"`
	AuthorID     string `json:"authorId" binding:"required,uuid"`
	CollectionID string `json:"collectionId" binding:"required,uuid"`
}

type UpdateArtworkRequest struct {
	Title        *string `json:"title,omitempty"`
	CreationYear *int    `json:"creationYear,omitempty"`
	Technic      *string `json:"technic,omitempty"`
	Material     *string `json:"material,omitempty"`
	Size         *string `json:"size,omitempty"`
	AuthorID     *string `json:"authorId,omitempty" validate:"omitempty,uuid"`
	CollectionID *string `json:"collectionId,omitempty" validate:"omitempty,uuid"`
}

var (
	ErrArtworkEmptyTitle        = errors.New("empty title")
	ErrArtworkTitleTooLong      = errors.New("title exceeds maximum length (255 chars)")
	ErrArtworkEmptyTechnic      = errors.New("empty technic")
	ErrArtworkTechnicTooLong    = errors.New("technic exceeds maximum length (100 chars)")
	ErrArtworkEmptyMaterial     = errors.New("empty material")
	ErrArtworkMaterialTooLong   = errors.New("material exceeds maximum length (100 chars)")
	ErrArtworkEmptySize         = errors.New("empty size")
	ErrArtworkSizeTooLong       = errors.New("size exceeds maximum length (50 chars)")
	ErrArtworkInvalidYear       = errors.New("invalid creation year")
	ErrArtworkInvalidAuthor     = errors.New("invalid author reference")
	ErrArtworkInvalidCollection = errors.New("invalid collection reference")
	ErrArtworkYearNotInRange    = errors.New("creation year not in author's lifetime")
)

func NewArtwork(
	id uuid.UUID,
	title string,
	technic string,
	material string,
	size string,
	creationYear int,
	author *Author,
	collection *Collection,
) (Artwork, error) {
	artwork := Artwork{
		id:           id,
		title:        strings.TrimSpace(title),
		technic:      strings.TrimSpace(technic),
		material:     strings.TrimSpace(material),
		size:         strings.TrimSpace(size),
		creationYear: creationYear,
		author:       author,
		collection:   collection,
	}

	if err := artwork.validate(); err != nil {
		return Artwork{}, err
	}

	if err := artwork.validateWithAuthor(); err != nil {
		return Artwork{}, err
	}

	return artwork, nil
}

func (a *Artwork) validate() error {
	switch {
	case a.title == "":
		return ErrArtworkEmptyTitle
	case len(a.title) > 255:
		return ErrArtworkTitleTooLong
	case a.technic == "":
		return ErrArtworkEmptyTechnic
	case len(a.technic) > 100:
		return ErrArtworkTechnicTooLong
	case a.material == "":
		return ErrArtworkEmptyMaterial
	case len(a.material) > 100:
		return ErrArtworkMaterialTooLong
	case a.size == "":
		return ErrArtworkEmptySize
	case len(a.size) > 50:
		return ErrArtworkSizeTooLong
	case a.creationYear < 0:
		return ErrArtworkInvalidYear
	case a.author == nil:
		return ErrArtworkInvalidAuthor
	case a.collection == nil:
		return ErrArtworkInvalidCollection
	}
	return nil
}

func (a *Artwork) validateWithAuthor() error {
	if a.author == nil {
		return ErrArtworkInvalidAuthor
	}

	birthYear := a.author.GetBirthYear()
	deathYear := a.author.GetDeathYear()

	if a.creationYear < birthYear || (deathYear > 0 && a.creationYear > deathYear) {
		return ErrArtworkYearNotInRange
	}

	return nil
}

func (a *Artwork) ToArtworkResponse() ArtworkResponse {
	return ArtworkResponse{
		ID:           a.id.String(),
		Title:        a.title,
		CreationYear: a.creationYear,
		Technic:      a.technic,
		Material:     a.material,
		Size:         a.size,
		Author:       a.GetAuthor().ToAuthorResponse(),
		Collection:   a.GetCollection().ToCollectionResponse(),
	}
}

// func ToArtworkModel(req CreateArtworkRequest) (Artwork, error) {
// 	authorID, err := uuid.Parse(req.AuthorID)
// 	if err != nil {
// 		return Artwork{}, err
// 	}

// 	collectionID, err := uuid.Parse(req.CollectionID)
// 	if err != nil {
// 		return Artwork{}, err
// 	}

// 	return Artwork{
// 		id:           uuid.New(),
// 		title:        req.Title,
// 		creationYear: req.CreationYear,
// 		technic:      req.Technic,
// 		material:     req.Material,
// 		size:         req.Size,
// 		author:       &Author{ID: authorID},
// 		collection:   &Collection{ID: collectionID},
// 	}, nil
// }

func (a *Artwork) GetID() uuid.UUID {
	return a.id
}

func (a *Artwork) GetTitle() string {
	return a.title
}

func (a *Artwork) GetCreationYear() int {
	return a.creationYear
}

func (a *Artwork) GetAuthor() *Author {
	return a.author
}

func (a *Artwork) GetCollection() *Collection {
	return a.collection
}

func (a *Artwork) GetSize() string {
	return a.size
}

func (a *Artwork) GetMaterial() string {
	return a.material
}

func (a *Artwork) GetTechnic() string {
	return a.technic
}
