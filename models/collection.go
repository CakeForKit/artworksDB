package models

import (
	"errors"

	"github.com/google/uuid"
)

type Collection struct {
	id    uuid.UUID
	title string
}

var ErrCollectionEmptyTitle = errors.New("empty title")

func NewCollection(title string) (Collection, error) {
	if title == "" {
		return Collection{}, ErrAuthorEmptyName
	}
	return Collection{
		id:    uuid.New(),
		title: title,
	}, nil
}
