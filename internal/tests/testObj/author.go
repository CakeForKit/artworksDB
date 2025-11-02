package testobj

import (
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
)

type AuthorMother interface {
	AuthorP(authorID uuid.UUID) *models.Author
	AuthorUpdateReq() models.AuthorUpdateReq
	AuthorBirthYearP(authorID uuid.UUID, birthYear int) *models.Author
}

func NewAuthorMother() AuthorMother {
	return &authorMother{}
}

type authorMother struct {
}

func (um *authorMother) AuthorP(authorID uuid.UUID) *models.Author {
	author, _ := models.NewAuthor(
		authorID,
		"test-author-name",
		1970,
		1990,
	)
	return &author
}

func (um *authorMother) AuthorUpdateReq() models.AuthorUpdateReq {
	return models.AuthorUpdateReq{
		Name:      "changed-test-author-name",
		BirthYear: 1800,
		DeathYear: 1880,
	}
}

func (um *authorMother) AuthorBirthYearP(authorID uuid.UUID, birthYear int) *models.Author {
	author, _ := models.NewAuthor(
		authorID,
		"test-author-name",
		birthYear,
		birthYear+50,
	)
	return &author
}
