package jsonreqresp

import "github.com/google/uuid"

type StatCollectionsResponse struct {
	ColID       uuid.UUID `json:"ColID"`
	ColTitle    string    `json:"ColTitle"`
	CntArtworks int       `json:"CntArtworks" `
}
