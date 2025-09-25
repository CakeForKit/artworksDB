package fixturesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func getAdminIDsOfEmployees(employees []*models.Employee) (adminIDs uuid.UUIDs) {
	adminsMap := make(map[uuid.UUID]bool, 0)
	for _, v := range employees {
		id := v.GetAdminID()
		if _, ok := adminsMap[id]; !ok {
			adminsMap[id] = true
			adminIDs = append(adminIDs, id)
		}
	}
	return
}

func AddTestEmployees(
	t provider.T, ctx context.Context, employees []*models.Employee,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
) {
	// fmt.Printf("AddTestEmployees:\n")
	// fmt.Println("Employees:")
	// for _, v := range employees {
	// 	fmt.Printf("(%v) %v\n\n", v.GetID(), v)
	// }
	// fmt.Print("\n\n")

	t.WithNewStep("Add Test Admins", func(sCtx provider.StepCtx) {
		adminIDs := getAdminIDsOfEmployees(employees)
		adminCreator := testobj.NewAdminMother()
		admins := make([]*models.Admin, 0, len(adminIDs))
		for _, id := range adminIDs {
			admins = append(admins, adminCreator.DefaultAdminP(id))
		}
		AddTestAdmin(t, ctx, adminRep, admins)
	})
	t.WithNewStep("Add Test Employees", func(sCtx provider.StepCtx) {
		for _, u := range employees {
			err := employeeRep.Add(ctx, u)
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestEmployees(
	t provider.T, ctx context.Context, employees []*models.Employee,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
) {
	// fmt.Printf("DelTestEmployees:\n")
	// fmt.Println("Employees:")
	// for _, v := range employees {
	// 	fmt.Printf("(%v) %v\n\n", v.GetID(), v)
	// }
	// fmt.Print("\n\n")
	t.WithNewStep("Delete Test Employees", func(sCtx provider.StepCtx) {
		for _, u := range employees {
			err := employeeRep.Delete(ctx, u.GetID())
			sCtx.Assert().NoError(err)
		}
	})
	t.WithNewStep("Delete Test Admins", func(sCtx provider.StepCtx) {
		adminIDs := getAdminIDsOfEmployees(employees)
		admins := make([]*models.Admin, 0, len(adminIDs))
		for _, id := range adminIDs {
			a, err := adminRep.GetByID(ctx, id)
			sCtx.Assert().NoError(err)
			admins = append(admins, a)
		}
		DelTestAdmin(t, ctx, adminRep, admins)
	})
}

func AssertEmployeesAreInRes(t provider.T, employees, resEmployees []*models.Employee) {
	t.WithNewStep("Check employees sre in the result", func(sCtx provider.StepCtx) {
		foundAll := make([]bool, len(employees))
		for _, ru := range resEmployees {
			for i, u := range employees {
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
