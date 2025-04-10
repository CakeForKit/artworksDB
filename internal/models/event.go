package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	id        uuid.UUID
	title     string
	dateBegin time.Time
	dateEnd   time.Time
	adress    string
	access    bool
	artworks  []*Artwork
}

var (
	ErrEventEmptyTitle      = errors.New("empty title")
	ErrEventInvalidDates    = errors.New("end date must be after start date")
	ErrEventEmptyAddress    = errors.New("empty address")
	ErrEventInvalidAccess   = errors.New("invalid access value")
	ErrArtworkAlreadyExists = errors.New("artwork already exists in event")
	ErrArtworkNotFound      = errors.New("artwork not found in event")
)

func NewEvent(title string, dateBegin time.Time, dateEnd time.Time,
	address string, access bool, artworks []*Artwork) (Event, error) {
	if title == "" {
		return Event{}, ErrEventEmptyTitle
	}
	if dateEnd.Before(dateBegin) {
		return Event{}, ErrEventInvalidDates
	}
	if address == "" {
		return Event{}, ErrEventEmptyAddress
	}
	if access != true && access != false {
		return Event{}, ErrEventInvalidAccess
	}

	return Event{
		id:        uuid.New(),
		title:     title,
		dateBegin: dateBegin,
		dateEnd:   dateEnd,
		adress:    address,
		access:    access,
		artworks:  artworks,
	}, nil
}

func (e *Event) TextAbout() string {
	return fmt.Sprintf(
		"Событие: %s:\nДата начала: %s\nДата окончания: %s\nАдрес: %s\nДоступ: %v\n",
		e.GetTitle(),
		e.GetDateBegin().Format(time.RFC3339),
		e.GetDateEnd().Format(time.RFC3339),
		e.GetAddress(),
		e.GetAccess(),
	)
}

// GetID возвращает уникальный идентификатор события
func (e *Event) GetID() uuid.UUID {
	return e.id
}

// GetTitle возвращает название события
func (e *Event) GetTitle() string {
	return e.title
}

// GetDateBegin возвращает дату и время начала события
func (e *Event) GetDateBegin() time.Time {
	return e.dateBegin
}

// GetDateEnd возвращает дату и время окончания события
func (e *Event) GetDateEnd() time.Time {
	return e.dateEnd
}

// GetAddress возвращает адрес проведения события
// Примечание: исправлено с "adress" на "address" для соответствия правильному написанию
func (e *Event) GetAddress() string {
	return e.adress
}

// GetAccess возвращает флаг доступа к событию (публичное/приватное)
func (e *Event) GetAccess() bool {
	return e.access
}

// GetArtworks возвращает список произведений искусства на событии
func (e *Event) GetArtworks() []*Artwork {
	return e.artworks
}

// AddArtwork добавляет произведение искусства к событию
func (e *Event) AddArtwork(artwork *Artwork) error {
	for _, a := range e.artworks {
		if a.GetID() == artwork.GetID() {
			return ErrArtworkAlreadyExists
		}
	}
	e.artworks = append(e.artworks, artwork)
	return nil
}

// RemoveArtwork удаляет произведение искусства из события
func (e *Event) RemoveArtwork(artworkID uuid.UUID) error {
	for i, a := range e.artworks {
		if a.GetID() == artworkID {
			e.artworks = append(e.artworks[:i], e.artworks[i+1:]...)
			return nil
		}
	}
	return ErrArtworkNotFound
}

// CheckArtwork проверяет наличие произведения искусства в событии
func (e *Event) CheckArtwork(artworkID uuid.UUID) bool {
	for _, a := range e.artworks {
		if a.GetID() == artworkID {
			return true
		}
	}
	return false
}

// ClearArtworks очищает список произведений искусства
func (e *Event) ClearArtworks() {
	e.artworks = []*Artwork{}
}

// ArtworksCount возвращает количество произведений искусства на событии
func (e *Event) ArtworksCount() int {
	return len(e.artworks)
}
