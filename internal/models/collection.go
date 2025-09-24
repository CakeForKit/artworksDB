package models

import (
	"errors"
	"strings"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type Collection struct {
	id    uuid.UUID
	title string
}

type CollectionUpdateReq struct {
	Title string `json:"title" example:"Louvre Museum Collection"`
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

func (c1 *Collection) validate() error {
	switch {
	case c1.title == "":
		return ErrCollectionEmptyTitle
	case len(c1.title) > 255:
		return ErrCollectionTitleTooLong
	}
	return nil
}

func (c1 *Collection) Equal(other interface{}) bool {
	if c1 == nil {
		return other == nil
	}

	c2, ok := other.(*Collection)
	if !ok {
		return false
	}
	if c2 == nil {
		return false
	}

	return c1.title == c2.title
}

func (c1 *Collection) ToCollectionResponse() jsonreqresp.CollectionResponse {
	return jsonreqresp.CollectionResponse{
		ID:    c1.id.String(),
		Title: c1.title,
	}
}

// func FromCollectionRequest(req jsonreqresp.CollectionAddRequest) (Collection, error) {
// 	var id uuid.UUID
// 	if req.ID == "" {
// 		id = uuid.New()
// 	} else {
// 		var err error
// 		id, err = uuid.Parse(req.ID)
// 		if err != nil {
// 			return Collection{}, fmt.Errorf("FromAuthorRequest: %w", err)
// 		}
// 	}
// 	return NewCollection(id, req.Title)
// }

func (c1 *Collection) GetID() uuid.UUID {
	return c1.id
}

func (c1 *Collection) GetTitle() string {
	return c1.title
}

func (c1 *Collection) Update(updateReq CollectionUpdateReq) error {
	copyC := *c1
	copyC.title = updateReq.Title
	if err := copyC.validate(); err != nil {
		return err
	}
	*c1 = copyC
	return nil
}
