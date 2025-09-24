package integration

import (
	"context"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type EmployeeRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx             context.Context
	employeeCreator testobj.EmployeeMother
	employeeRep     employeerep.EmployeeRep
}

func TestEmployeeRepSuite(t *testing.T) {
	suite.RunSuite(t, new(EmployeeRepSuite))
}

func (s *EmployeeRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.employeeCreator = testobj.NewEmployeeMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.employeeRep, err = employeerep.NewEmployeeRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *EmployeeRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "employee")
}

func (s *EmployeeRepSuite) AfterAll(t provider.T) {
	if s.employeeRep != nil {
		s.employeeRep.Close()
	}
}

// func (s *EmployeeRepSuite) TestEmployeeRep_GetAll(t provider.T) {
// 	t.Parallel()
// 	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
// 		employees := []*models.Employee{
// 			s.employeeCreator.EmployeeP(uuid.New()),
// 			s.employeeCreator.EmployeeP(uuid.New()),
// 		}
// 		fixturesrep.AddTestEmployees(t, s.ctx, employees, s.employeeRep)
// 		defer fixturesrep.DelTestEmployees(t, s.ctx, employees, s.employeeRep)

// 		// ACT
// 		resEmployees, err := s.employeeRep.GetAll(s.ctx)

// 		sCtx.Require().NoError(err)
// 		fixturesrep.AssertEmployeesAreInRes(t, employees, resEmployees)
// 	})
// }
