package models

import (
	"errors"
	"fmt"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
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

func (s *StatCollections) ToResponse() jsonreqresp.StatCollectionsResponse {
	return jsonreqresp.StatCollectionsResponse{
		ColID:       s.colID,
		ColTitle:    s.colTitle,
		CntArtworks: s.cntArtworks,
	}
}
