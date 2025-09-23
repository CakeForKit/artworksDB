package fixturesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func AddTestCollections(
	t provider.T, ctx context.Context, collections []*models.Collection,
	collectionRep collectionrep.CollectionRep,
) {
	t.WithNewStep("Add Test Collections", func(sCtx provider.StepCtx) {
		for _, u := range collections {
			err := collectionRep.Add(ctx, u)
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestCollections(
	t provider.T, ctx context.Context, collections []*models.Collection,
	collectionRep collectionrep.CollectionRep,
) {
	t.WithNewStep("Del Test Collections", func(sCtx provider.StepCtx) {
		for _, u := range collections {
			err := collectionRep.Delete(ctx, u.GetID())
			sCtx.Require().NoError(err)
		}
	})
}
