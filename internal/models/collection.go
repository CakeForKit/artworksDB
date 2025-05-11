package models

import (
	"errors"
	"fmt"
	"strings"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
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

func (c *Collection) ToCollectionResponse() jsonreqresp.CollectionResponse {
	return jsonreqresp.CollectionResponse{
		ID:    c.id.String(),
		Title: c.title,
	}
}

func FromCollectionRequest(req jsonreqresp.CollectionRequest) (Collection, error) {
	var id uuid.UUID
	if req.ID == "" {
		id = uuid.New()
	} else {
		var err error
		id, err = uuid.Parse(req.ID)
		if err != nil {
			return Collection{}, fmt.Errorf("FromAuthorRequest: %w", err)
		}
	}
	return NewCollection(id, req.Title)
}

func (c *Collection) GetID() uuid.UUID {
	return c.id
}

func (c *Collection) GetTitle() string {
	return c.title
}
