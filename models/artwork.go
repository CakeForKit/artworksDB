package models

import (
	"errors"

	"github.com/google/uuid"
)

type Artwork struct {
	id           uuid.UUID
	title        string
	creationYear int
	author       *Author
	collection   *Collection
	size         string
	material     string
	technic      string
}

var (
	ErrArtworkEmptyTitle           = errors.New("empty title")
	ErrArtworkInvalidYear          = errors.New("invalid year")
	ErrArtworkInvalidAuthor        = errors.New("invalid author")
	ErrArtworkInvalidCollection    = errors.New("invalid collection")
	ErrArtworkInvalidSize          = errors.New("empty size")
	ErrArtworkInvalidMaterial      = errors.New("empty material")
	ErrArtworkInvalidTechnic       = errors.New("empty technic")
	ErrArtworkCreationYearToAuthor = errors.New("invalid creation year to author live")
)

func NewArtwork(title string, creationYear int, author *Author, collection *Collection,
	size string, material string, technic string) (Artwork, error) {

	if title == "" {
		return Artwork{}, ErrArtworkEmptyTitle
	} else if creationYear < 0 {
		return Artwork{}, ErrArtworkInvalidYear
	} else if author == nil {
		return Artwork{}, ErrArtworkInvalidAuthor
	} else if !(author.birthYear < creationYear && creationYear <= author.deathYear) {
		return Artwork{}, ErrArtworkCreationYearToAuthor
	} else if collection == nil {
		return Artwork{}, ErrArtworkInvalidCollection
	} else if size == "" {
		return Artwork{}, ErrArtworkInvalidSize
	} else if material == "" {
		return Artwork{}, ErrArtworkInvalidMaterial
	} else if technic == "" {
		return Artwork{}, ErrArtworkInvalidTechnic
	}
	return Artwork{
		id:           uuid.New(),
		title:        title,
		creationYear: creationYear,
		author:       author,
		collection:   collection,
		size:         size,
		material:     material,
		technic:      technic,
	}, nil
}

func (a *Artwork) GetID() uuid.UUID {
	return a.id
}
