package fixturesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func AddTestArtworks(
	t provider.T, ctx context.Context, artworks []*models.Artwork,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {

	t.WithNewStep("Add Test Authors", func(sCtx provider.StepCtx) {
		authors := make(map[uuid.UUID]*models.Author, 0)
		for _, v := range artworks {
			if _, ok := authors[v.GetAuthor().GetID()]; !ok {
				authors[v.GetAuthor().GetID()] = v.GetAuthor()
				err := authorRep.Add(ctx, v.GetAuthor())
				sCtx.Require().NoError(err)
			}
		}
	})
	t.WithNewStep("Add Test Collections", func(sCtx provider.StepCtx) {
		cols := make(map[uuid.UUID]*models.Collection, 0)
		for _, v := range artworks {
			if _, ok := cols[v.GetCollection().GetID()]; !ok {
				cols[v.GetCollection().GetID()] = v.GetCollection()
				err := collectionRep.Add(ctx, v.GetCollection())
				sCtx.Require().NoError(err)
			}
		}
	})
	t.WithNewStep("Add Test Artworks", func(sCtx provider.StepCtx) {
		for _, v := range artworks {
			err := artworkRep.Add(ctx, v)
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestArtworks(
	t provider.T, ctx context.Context, artworks []*models.Artwork,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {
	t.WithNewStep("Del Test Artworks", func(sCtx provider.StepCtx) {
		for _, v := range artworks {
			err := artworkRep.Delete(ctx, v.GetID())
			sCtx.Require().NoError(err)
		}
	})
	t.WithNewStep("Del Test Authors", func(sCtx provider.StepCtx) {
		authors := make(map[uuid.UUID]*models.Author, 0)
		for _, v := range artworks {
			if _, ok := authors[v.GetAuthor().GetID()]; !ok {
				authors[v.GetAuthor().GetID()] = v.GetAuthor()
				err := authorRep.Delete(ctx, v.GetAuthor().GetID())
				sCtx.Require().NoError(err)
			}
		}
	})
	t.WithNewStep("Del Test Collections", func(sCtx provider.StepCtx) {
		cols := make(map[uuid.UUID]*models.Collection, 0)
		for _, v := range artworks {
			if _, ok := cols[v.GetCollection().GetID()]; !ok {
				cols[v.GetCollection().GetID()] = v.GetCollection()
				err := collectionRep.Delete(ctx, v.GetCollection().GetID())
				sCtx.Require().NoError(err)
			}
		}
	})

}
