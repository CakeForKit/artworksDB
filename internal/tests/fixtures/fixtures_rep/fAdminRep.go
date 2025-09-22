package fixturesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func AddTestAdmin(t provider.T, ctx context.Context, adminRep adminrep.AdminRep, admins []*models.Admin) {
	t.WithNewStep("Add Test Admins", func(sCtx provider.StepCtx) {
		for _, u := range admins {
			err := adminRep.Add(ctx, u)
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestAdmin(t provider.T, ctx context.Context, adminRep adminrep.AdminRep, admins []*models.Admin) {
	t.WithNewStep("Delete Test Admins", func(sCtx provider.StepCtx) {
		for _, u := range admins {
			err := adminRep.Delete(ctx, u.GetID())
			sCtx.Assert().NoError(err)
		}
	})
}

// func AssertAdminsAreInRes(t provider.T, admins, resAdmins []*models.Admin) {
// 	t.WithNewStep("Check admins sre in the result", func(sCtx provider.StepCtx) {
// 		foundAll := make([]bool, len(admins))
// 		for _, ru := range resAdmins {
// 			for i, u := range admins {
// 				if models.CmpAdmins(ru, u) {
// 					foundAll[i] = true
// 				}
// 			}
// 		}
// 		for _, v := range foundAll {
// 			sCtx.Assert().True(v)
// 		}
// 	})
// }
