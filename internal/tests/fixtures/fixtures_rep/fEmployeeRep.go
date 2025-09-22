package fixturesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func AddTestEmployee(t provider.T, ctx context.Context, employeeRep employeerep.EmployeeRep, employees []*models.Employee) {
	t.WithNewStep("Add Test Employees", func(sCtx provider.StepCtx) {
		for _, u := range employees {
			err := employeeRep.Add(ctx, u)
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestEmployee(t provider.T, ctx context.Context, employeeRep employeerep.EmployeeRep, employees []*models.Employee) {
	t.WithNewStep("Delete Test Employees", func(sCtx provider.StepCtx) {
		for _, u := range employees {
			err := employeeRep.Delete(ctx, u.GetID())
			sCtx.Assert().NoError(err)
		}
	})
}

// func AssertEmployeesAreInRes(t provider.T, employees, resEmployees []*models.Employee) {
// 	t.WithNewStep("Check employees sre in the result", func(sCtx provider.StepCtx) {
// 		foundAll := make([]bool, len(employees))
// 		for _, ru := range resEmployees {
// 			for i, u := range employees {
// 				if models.CmpEmployees(ru, u) {
// 					foundAll[i] = true
// 				}
// 			}
// 		}
// 		for _, v := range foundAll {
// 			sCtx.Assert().True(v)
// 		}
// 	})
// }
