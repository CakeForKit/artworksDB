package fixturesrep

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func AddTestUser(t provider.T, ctx context.Context, userRep userrep.UserRep, users []*models.User) {
	t.WithNewStep("Add Test Users", func(sCtx provider.StepCtx) {
		for _, u := range users {
			err := userRep.Add(ctx, u)
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestUser(t provider.T, ctx context.Context, userRep userrep.UserRep, users []*models.User) {
	t.WithNewStep("Delete Test Users", func(sCtx provider.StepCtx) {
		for _, u := range users {
			err := userRep.Delete(ctx, u.GetID())
			sCtx.Assert().NoError(err)
		}
	})
}

func AssertUsersAreInRes(t provider.T, users, resUsers []*models.User) {
	t.WithNewStep("Check users sre in the result", func(sCtx provider.StepCtx) {
		foundAll := make([]bool, len(users))
		for _, ru := range resUsers {
			for i, u := range users {
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
