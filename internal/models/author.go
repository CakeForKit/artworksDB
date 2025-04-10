package models

import (
	"errors"

	"github.com/google/uuid"
)

type Author struct {
	id        uuid.UUID
	name      string
	birthYear int
	deathYear int
}

var ErrAuthorEmptyName = errors.New("empty name")
var ErrAuthorInvalidYear = errors.New("invalid year")

func NewAuthor(name string, birthYear int, deathYear int) (Author, error) {
	if name == "" {
		return Author{}, ErrAuthorEmptyName
	} else if birthYear < 0 || deathYear < 0 || deathYear <= birthYear {
		return Author{}, ErrAuthorInvalidYear
	}
	return Author{
		id:        uuid.New(),
		name:      name,
		birthYear: birthYear,
		deathYear: deathYear,
	}, nil
}

func (auth *Author) GetID() uuid.UUID {
	return auth.id
}

func (auth *Author) SetID(id uuid.UUID) {
	auth.id = id
}

func (auth *Author) GetName() string {
	return auth.name
}

func (auth *Author) SetName(name string) {
	auth.name = name
}
