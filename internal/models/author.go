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

type AuthorResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"name" example:"Leonardo da Vinci"`
	BirthYear int    `json:"birthYear" example:"1452"`
	DeathYear int    `json:"deathYear" example:"1519"`
}

type AuthorRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=100"`                      // Обязательное, 2-100 символов
	BirthYear int    `json:"birthYear" binding:"required,gte=1000"`                      // Обязательное, >= 1000
	DeathYear *int   `json:"deathYear,omitempty" binding:"omitempty,gtefield=BirthYear"` // Опциональное, >= BirthYear
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

func (a *Author) ToAuthorResponse() AuthorResponse {
	return AuthorResponse{
		ID:        a.id.String(),
		Name:      a.name,
		BirthYear: a.birthYear,
		DeathYear: a.deathYear,
	}
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
