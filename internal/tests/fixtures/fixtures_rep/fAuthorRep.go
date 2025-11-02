package fixturesrep

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func AddTestAuthors(
	t provider.T, ctx context.Context, authors []*models.Author,
	authorRep authorrep.AuthorRep,
) {
	t.WithNewStep("Add Test Authors", func(sCtx provider.StepCtx) {
		for _, u := range authors {
			err := authorRep.Add(ctx, u)
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestAuthors(
	t provider.T, ctx context.Context, authors []*models.Author,
	authorRep authorrep.AuthorRep,
) {
	t.WithNewStep("Del Test Authors", func(sCtx provider.StepCtx) {
		for _, u := range authors {
			err := authorRep.Delete(ctx, u.GetID())
			sCtx.Require().NoError(err)
		}
	})
}

func AssertAuthorsAreInRes(t provider.T, authors, resAuthors []*models.Author) {
	t.WithNewStep("Check authors sre in the result", func(sCtx provider.StepCtx) {
		foundAll := make([]bool, len(authors))
		for _, ru := range resAuthors {
			for i, u := range authors {
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
