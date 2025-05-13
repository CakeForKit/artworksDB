package models

import (
	"errors"
	"fmt"
	"strings"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
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

var (
	ErrValidateArtwork          = errors.New("invalid model")
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

	// validateWithAuthor
	birthYear := a.author.GetBirthYear()
	deathYear := a.author.GetDeathYear()
	fmt.Printf("deathYear = %d\n", deathYear)
	if a.creationYear < birthYear || (deathYear > 0 && a.creationYear > deathYear) {
		return ErrArtworkYearNotInRange
	}
	return nil
}

func (a *Artwork) ToArtworkResponse() jsonreqresp.ArtworkResponse {
	return jsonreqresp.ArtworkResponse{
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

func (a *Artwork) Update(updateReq jsonreqresp.ArtworkUpdate) error {
	copyA := *a
	copyA.title = updateReq.Title
	copyA.creationYear = updateReq.CreationYear
	copyA.technic = updateReq.Technic
	copyA.material = updateReq.Material
	copyA.size = updateReq.Size

	if err := copyA.validate(); err != nil {
		return err
	}

	*a = copyA
	return nil
}
