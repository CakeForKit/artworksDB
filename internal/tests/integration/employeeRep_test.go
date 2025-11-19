package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	"github.com/CakeForKit/artworksDB.git/internal/tests/fixtures"
	fixturesrep "github.com/CakeForKit/artworksDB.git/internal/tests/fixtures/fixtures_rep"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type EmployeeRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx             context.Context
	employeeCreator testobj.EmployeeMother
	adminCreator    testobj.AdminMother
	employeeRep     employeerep.EmployeeRep
	adminRep        adminrep.AdminRep
}

func TestEmployeeRepSuite(t *testing.T) {
	suite.RunSuite(t, new(EmployeeRepSuite))
}

func (s *EmployeeRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.employeeCreator = testobj.NewEmployeeMother()
	s.adminCreator = testobj.NewAdminMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error
		s.employeeRep, err = employeerep.NewEmployeeRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.adminRep, err = adminrep.NewAdminRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
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

func (s *EmployeeRepSuite) TestEmployeeRep_GetAll(t provider.T) {
	t.Parallel()
	t.Run("Success", func(t provider.T) {
		employees := []*models.Employee{
			s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
			s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
		}
		fixturesrep.AddTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)

		// ACT
		resEmployees, err := s.employeeRep.GetAll(s.ctx)

		t.Require().NoError(err)
		fixturesrep.AssertEmployeesAreInRes(t, employees, resEmployees)
	})
}

func (s *EmployeeRepSuite) TestEmployeeRep_GetByID(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		employees := []*models.Employee{
			s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
		}
		fixturesrep.AddTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)

		// ACT
		resEmployee, err := s.employeeRep.GetByID(s.ctx, employees[0].GetID())

		t.Require().NoError(err)
		t.Require().True(resEmployee.Equal(employees[0]))
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.employeeRep.GetByID(s.ctx, uuid.New())

		t.Require().ErrorIs(err, employeerep.ErrEmployeeNotFound)
	})
}

func (s *EmployeeRepSuite) TestEmployeeRep_GetByLogin(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		employees := []*models.Employee{
			s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
		}
		fixturesrep.AddTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)

		// ACT
		resEmployee, err := s.employeeRep.GetByLogin(s.ctx, employees[0].GetLogin())

		t.Require().NoError(err)
		t.Require().True(resEmployee.Equal(employees[0]))
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.employeeRep.GetByLogin(s.ctx, "nonexistent_login")

		t.Require().ErrorIs(err, employeerep.ErrEmployeeNotFound)
	})
}

func (s *EmployeeRepSuite) TestEmployeeRep_Add(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		admin := s.adminCreator.DefaultAdminP(uuid.New())
		employee := s.employeeCreator.DefaultEmployeeP(uuid.New(), admin.GetID())
		err := s.adminRep.Add(s.ctx, admin)
		t.Require().NoError(err)
		defer fixturesrep.DelTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)

		// ACT
		err = s.employeeRep.Add(s.ctx, employee)

		t.Require().NoError(err)

		// Verify employee was added
		resEmployee, err := s.employeeRep.GetByID(s.ctx, employee.GetID())
		t.Require().NoError(err)
		t.Require().True(resEmployee.Equal(employee))
	})

	t.Run("Duplicate login error", func(t provider.T) {
		employee := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())
		fixturesrep.AddTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)

		// ACT - try to add employee with same login
		duplicateEmployee := s.employeeCreator.EmployeeWithLoginP(uuid.New(), uuid.New(), employee.GetLogin())
		err := s.employeeRep.Add(s.ctx, duplicateEmployee)

		t.Require().Error(err)
	})

	t.Run("Duplicate ID error", func(t provider.T) {
		employee := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())
		fixturesrep.AddTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)

		// ACT - try to add employee with same ID
		err := s.employeeRep.Add(s.ctx, employee)

		t.Require().Error(err)
	})
}

func (s *EmployeeRepSuite) TestEmployeeRep_Delete(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		employee := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())
		fixturesrep.AddTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)

		// ACT
		err := s.employeeRep.Delete(s.ctx, employee.GetID())

		t.Require().NoError(err)

		// Verify employee was deleted
		_, err = s.employeeRep.GetByID(s.ctx, employee.GetID())
		t.Require().ErrorIs(err, employeerep.ErrEmployeeNotFound)
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.employeeRep.Delete(s.ctx, uuid.New())

		t.Require().Error(err)
	})
}

func (s *EmployeeRepSuite) TestEmployeeRep_Update(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		employee := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())
		fixturesrep.AddTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)

		newEmployee, err := models.NewEmployee(
			employee.GetID(),
			"Updated Employee Name",
			"updated_login_"+uuid.NewString(),
			employee.GetHashedPassword(),
			employee.GetCreatedAt(),
			true,
			employee.GetAdminID(),
		)
		t.Require().NoError(err)
		funcUpdate := func(e *models.Employee) (*models.Employee, error) {
			return &newEmployee, nil
		}

		// ACT
		updatedEmployee, err := s.employeeRep.Update(s.ctx, employee.GetID(), funcUpdate)

		t.Require().NoError(err)
		t.Require().True(newEmployee.Equal(updatedEmployee))

		// Verify changes persisted
		dbEmployee, err := s.employeeRep.GetByID(s.ctx, employee.GetID())
		t.Require().NoError(err)
		t.Require().True(newEmployee.Equal(dbEmployee))
	})

	t.Run("Not found", func(t provider.T) {
		funcUpdate := func(e *models.Employee) (*models.Employee, error) {
			return e, nil
		}

		// ACT
		_, err := s.employeeRep.Update(s.ctx, uuid.New(), funcUpdate)

		t.Require().ErrorIs(err, employeerep.ErrEmployeeNotFound)
	})

	t.Run("Update function returns error", func(t provider.T) {
		employee := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())
		fixturesrep.AddTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)

		expectedErr := errors.New("update function error")
		funcUpdate := func(e *models.Employee) (*models.Employee, error) {
			return nil, expectedErr
		}

		// ACT
		_, err := s.employeeRep.Update(s.ctx, employee.GetID(), funcUpdate)

		t.Require().ErrorIs(err, expectedErr)
	})

	t.Run("Duplicate login error", func(t provider.T) {
		employee1 := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())
		employee2 := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())

		employees := []*models.Employee{employee1, employee2}
		fixturesrep.AddTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, employees, s.employeeRep, s.adminRep)

		newEmployee, err := models.NewEmployee(
			employee1.GetID(),
			employee1.GetUsername(),
			employee2.GetLogin(),
			employee1.GetHashedPassword(),
			employee1.GetCreatedAt(),
			employee1.IsValid(),
			employee1.GetAdminID(),
		)
		t.Require().NoError(err)
		funcUpdate := func(e *models.Employee) (*models.Employee, error) {
			return &newEmployee, nil
		}

		// ACT
		_, err = s.employeeRep.Update(s.ctx, employee1.GetID(), funcUpdate)

		t.Require().Error(err)
	})
}
