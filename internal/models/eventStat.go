package models

import (
	"errors"
	"fmt"

	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/google/uuid"
)

type StatCollections struct {
	colID       uuid.UUID
	colTitle    string
	cntArtworks int
}

var (
	ErrValidateStatCollections = errors.New("invalid model StatCollections")
	ErrCntArtworks             = errors.New("invalid cntArtoworks")
)

func NewStatCollections(colID uuid.UUID, colTitle string, cntArtworks int) (StatCollections, error) {
	colStat := StatCollections{
		colID:       colID,
		colTitle:    colTitle,
		cntArtworks: cntArtworks,
	}
	if err := colStat.validate(); err != nil {
		return StatCollections{}, fmt.Errorf("%w: %w", ErrValidateStatCollections, err)
	}
	return colStat, nil
}

func (s *StatCollections) validate() error {
	if s.cntArtworks <= 0 {
		return ErrCntArtworks
	}
	return nil
}
func (s1 *StatCollections) Equals(other interface{}) bool {
	if s1 == nil {
		return other == nil
	}

	s2, ok := other.(*StatCollections)
	if !ok {
		return false
	}
	if s2 == nil {
		return false
	}
	return s1.colID == s2.colID &&
		s1.colTitle == s2.colTitle &&
		s1.cntArtworks == s2.cntArtworks
}

func (s *StatCollections) ToResponse() jsonreqresp.StatCollectionsResponse {
	return jsonreqresp.StatCollectionsResponse{
		ColID:       s.colID,
		ColTitle:    s.colTitle,
		CntArtworks: s.cntArtworks,
	}
}

func (s *StatCollections) CollectionID() uuid.UUID {
	return s.colID
}

func (s *StatCollections) CollectionTitle() string {
	return s.colTitle
}

func (s *StatCollections) ArtworksCount() int {
	return s.cntArtworks
}
