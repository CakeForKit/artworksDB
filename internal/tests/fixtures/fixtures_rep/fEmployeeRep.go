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

func getAdminsOfEmployees(employees []*models.Employee) (admins []*models.Admin) {
	adminCreator := testobj.NewAdminMother()
	adminsMap := make(map[uuid.UUID]*models.Admin, 0)
	for _, v := range employees {
		if _, ok := adminsMap[v.GetAdminID()]; !ok {
			a := adminCreator.DefaultAdminP(v.GetAdminID())
			admins = append(admins, a)
			adminsMap[v.GetAdminID()] = a
		}
	}
	return
}

func AddTestEmployees(
	t provider.T, ctx context.Context, employees []*models.Employee,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
) {
	admins := getAdminsOfEmployees(employees)
	AddTestAdmin(t, ctx, adminRep, admins)

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
	t.WithNewStep("Delete Test Employees", func(sCtx provider.StepCtx) {
		for _, u := range employees {
			err := employeeRep.Delete(ctx, u.GetID())
			sCtx.Assert().NoError(err)
		}
	})
	admins := getAdminsOfEmployees(employees)
	DelTestAdmin(t, ctx, adminRep, admins)
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
