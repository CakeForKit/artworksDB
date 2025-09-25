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

func getAuthorsOfArtwork(artworks []*models.Artwork) (authors []*models.Author) {
	authorsMap := make(map[uuid.UUID]*models.Author, 0)
	for _, a := range artworks {
		author := a.GetAuthor()
		if _, ok := authorsMap[author.GetID()]; !ok {
			authorsMap[author.GetID()] = author
			authors = append(authors, author)
		}
	}
	return
}

func getCollectionsOfArtwork(artworks []*models.Artwork) (collections []*models.Collection) {
	collectionsMap := make(map[uuid.UUID]*models.Collection, 0)
	for _, a := range artworks {
		collection := a.GetCollection()
		if _, ok := collectionsMap[collection.GetID()]; !ok {
			collectionsMap[collection.GetID()] = collection
			collections = append(collections, collection)
		}
	}
	return
}

func AddTestArtworks(
	t provider.T, ctx context.Context, artworks []*models.Artwork,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {
	// fmt.Println("AddTestArtworks:")
	// fmt.Println("artworks:")
	// for _, v := range artworks {
	// 	fmt.Printf("%v -> %v\n\n", v, v.GetAuthor())
	// }
	// fmt.Print("\n\n")

	t.WithNewStep("Add Test Authors", func(sCtx provider.StepCtx) {
		authors := getAuthorsOfArtwork(artworks)
		AddTestAuthors(t, ctx, authors, authorRep)
	})
	t.WithNewStep("Add Test Collections", func(sCtx provider.StepCtx) {
		collections := getCollectionsOfArtwork(artworks)
		AddTestCollections(t, ctx, collections, collectionRep)
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
	// fmt.Println("DelTestArtworks:")
	// fmt.Println("artworks:")
	// for _, v := range artworks {
	// 	fmt.Printf("%v -> %v\n\n", v, v.GetAuthor())
	// }
	// fmt.Print("\n\n")

	t.WithNewStep("Del Test Artworks", func(sCtx provider.StepCtx) {
		for _, v := range artworks {
			err := artworkRep.Delete(ctx, v.GetID())
			sCtx.Assert().NoError(err)
		}
	})
	t.WithNewStep("Del Test Authors", func(sCtx provider.StepCtx) {
		authors := getAuthorsOfArtwork(artworks)
		DelTestAuthors(t, ctx, authors, authorRep)
	})
	t.WithNewStep("Del Test Collections", func(sCtx provider.StepCtx) {
		collections := getCollectionsOfArtwork(artworks)
		DelTestCollections(t, ctx, collections, collectionRep)
	})

}

func AssertArtworksAreInRes(t provider.T, artworks, resArtworks []*models.Artwork) {
	t.WithNewStep("Check artworks sre in the result", func(sCtx provider.StepCtx) {
		foundAll := make([]bool, len(artworks))
		for _, ru := range resArtworks {
			for i, u := range artworks {
				if ru.Equal(u) {
					foundAll[i] = true
				}
			}
		}
		for _, v := range foundAll {
			sCtx.Assert().True(v)
		}
	})
}
