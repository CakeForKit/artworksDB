package jsonreqresp

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ArtworkFilter struct {
	Title      string
	AuthorName string
	Collection string
	EventID    uuid.UUID
}

var (
	ErrEventFilterDate = errors.New("error format date for EventFilter")
)

type EventFilter struct {
	Title     string
	DateBegin time.Time
	DateEnd   time.Time
	CanVisit  string
	Valid     string
}

const (
	TitleSortFieldArtwork        = "title"
	AuthorNameSortFieldArtwork   = "author_name"
	CreationYearSortFieldArtwork = "creationYear"
	ASCDirection                 = "ASC"
	DESCDirection                = "DESC"
)

type ArtworkSortOps struct {
	Field     string `json:"field,omitempty" binding:"omitempty,oneof=title author_name creationYear" example:""` // обязательное, одно из значений
	Direction string `json:"direction,omitempty" binding:"omitempty,oneof=ASC DESC" example:""`                   // обязательное, только ASC или DESC
}
