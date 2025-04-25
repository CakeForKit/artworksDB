package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Collection struct {
	id    uuid.UUID
	title string
}

var (
	ErrCollectionEmptyTitle   = errors.New("empty title")
	ErrCollectionTitleTooLong = errors.New("title exceeds maximum length (255 chars)")
)

func NewCollection(id uuid.UUID, title string) (Collection, error) {
	collection := Collection{
		id:    id,
		title: strings.TrimSpace(title),
	}

	if err := collection.validate(); err != nil {
		return Collection{}, err
	}

	return collection, nil
}

func (c *Collection) validate() error {
	switch {
	case c.title == "":
		return ErrCollectionEmptyTitle
	case len(c.title) > 255:
		return ErrCollectionTitleTooLong
	}
	return nil
}

func (c *Collection) GetID() uuid.UUID {
	return c.id
}

func (c *Collection) GetTitle() string {
	return c.title
}
