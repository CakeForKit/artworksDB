package authorserv

import (
	"context"
	"errors"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"github.com/google/uuid"
)

type AuthorServ interface {
	GetAll(ctx context.Context) ([]*models.Author, error)
	Add(ctx context.Context, author *models.Author) error
	Update(ctx context.Context, idAuthor uuid.UUID, updateReq models.AuthorUpdateReq) error
	Delete(ctx context.Context, idAuthor uuid.UUID) error
}

var (
	ErrHasLinkedArtworks = errors.New("author has linked artworks")
)

func NewAuthorServ(authorRep authorrep.AuthorRep) AuthorServ {
	return &authorServ{authorRep: authorRep}
}

type authorServ struct {
	authorRep authorrep.AuthorRep
}

func (s *authorServ) GetAll(ctx context.Context) ([]*models.Author, error) {
	return s.authorRep.GetAll(ctx)
}

func (s *authorServ) Add(ctx context.Context, author *models.Author) error {
	return s.authorRep.Add(ctx, author)
}

func (s *authorServ) Update(ctx context.Context, idAuthor uuid.UUID, updateReq models.AuthorUpdateReq) error {
	return s.authorRep.Update(ctx, idAuthor, func(a *models.Author) (*models.Author, error) {
		err := a.Update(updateReq)
		return a, err
	})
}

func (s *authorServ) Delete(ctx context.Context, idAuthor uuid.UUID) error {
	has, err := s.authorRep.HasArtworks(ctx, idAuthor)
	if err != nil {
		return fmt.Errorf("authorServ.Delete %v", err)
	} else if has {
		return ErrHasLinkedArtworks
	}
	return s.authorRep.Delete(ctx, idAuthor)
}
