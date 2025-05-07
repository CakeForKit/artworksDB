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

type CollectionResponse struct {
	ID    string `json:"id" example:"aa1e8400-e29b-41d4-a716-446655441111"`
	Title string `json:"title" example:"Louvre Museum Collection"`
}

type CollectionRequest struct {
	Title string `json:"title" binding:"required,min=2,max=255"`
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

func (c *Collection) ToCollectionResponse() CollectionResponse {
	return CollectionResponse{
		ID:    c.id.String(),
		Title: c.title,
	}
}

func (c *Collection) GetID() uuid.UUID {
	return c.id
}

func (c *Collection) GetTitle() string {
	return c.title
}
