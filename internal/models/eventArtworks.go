package models

import (
	"errors"

	"github.com/google/uuid"
)

type EventArtworks struct {
	event    *Event
	artworks []*Artwork
}

var (
	ErrEventArtworksNilEvent   = errors.New("event cannot be nil")
	ErrEventArtworksNilArtwork = errors.New("artwork cannot be nil")
	ErrArtworkAlreadyInEvent   = errors.New("artwork already exists in event")
	ErrArtworkNotInEvent       = errors.New("artwork not found in event")
)

func NewEventArtworks(event *Event) (*EventArtworks, error) {
	if event == nil {
		return nil, ErrEventArtworksNilEvent
	}
	return &EventArtworks{
		event:    event,
		artworks: make([]*Artwork, 0),
	}, nil
}

func (ea *EventArtworks) GetEvent() *Event {
	return ea.event
}

func (ea *EventArtworks) GetArtworks() []*Artwork {
	return ea.artworks
}

func (ea *EventArtworks) AddArtwork(artwork *Artwork) error {
	if artwork == nil {
		return ErrEventArtworksNilArtwork
	}

	for _, a := range ea.artworks {
		if a.GetID() == artwork.GetID() {
			return ErrArtworkAlreadyInEvent
		}
	}
	ea.artworks = append(ea.artworks, artwork)
	return nil
}

func (ea *EventArtworks) RemoveArtwork(artworkID uuid.UUID) error {
	for i, a := range ea.artworks {
		if a.GetID() == artworkID {
			ea.artworks[i] = nil
			ea.artworks = append(ea.artworks[:i], ea.artworks[i+1:]...)
			return nil
		}
	}
	return ErrArtworkNotInEvent
}

func (ea *EventArtworks) HasArtwork(artworkID uuid.UUID) bool {
	for _, a := range ea.artworks {
		if a.GetID() == artworkID {
			return true
		}
	}
	return false
}

func (ea *EventArtworks) ClearArtworks() {
	ea.artworks = make([]*Artwork, 0)
}

func (ea *EventArtworks) CountArtworks() int {
	return len(ea.artworks)
}

func (ea *EventArtworks) IsEmpty() bool {
	return len(ea.artworks) == 0
}
