package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Author struct {
	id        uuid.UUID
	name      string
	birthYear int
	deathYear int
}

var (
	ErrAuthorEmptyName        = errors.New("empty name")
	ErrAuthorNameTooLong      = errors.New("name exceeds maximum length (100 chars)")
	ErrAuthorInvalidBirthYear = errors.New("invalid birth year")
	ErrAuthorInvalidDeathYear = errors.New("invalid death year")
	ErrAuthorBirthAfterDeath  = errors.New("birth year cannot be after death year")
	ErrAuthorLivingAuthor     = errors.New("for living authors, death year should be 0")
)

func NewAuthor(id uuid.UUID, name string, birthYear int, deathYear int) (Author, error) {
	author := Author{
		id:        id,
		name:      strings.TrimSpace(name),
		birthYear: birthYear,
		deathYear: deathYear,
	}
	if err := author.validate(); err != nil {
		return Author{}, err
	}
	return author, nil
}

func (a *Author) validate() error {
	switch {
	case a.name == "":
		return ErrAuthorEmptyName
	case len(a.name) > 100:
		return ErrAuthorNameTooLong
	case a.birthYear <= 0:
		return ErrAuthorInvalidBirthYear
	case a.deathYear < 0:
		return ErrAuthorInvalidDeathYear
	case a.deathYear > 0 && a.birthYear > a.deathYear:
		return ErrAuthorBirthAfterDeath
	}
	return nil
}

func (auth *Author) GetID() uuid.UUID {
	return auth.id
}

func (auth *Author) GetName() string {
	return auth.name
}

func (a *Author) GetBirthYear() int {
	return a.birthYear
}

func (a *Author) GetDeathYear() int {
	return a.deathYear
}
