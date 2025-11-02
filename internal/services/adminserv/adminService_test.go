package adminserv_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/adminserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

func TestAdminService_GetAllUsers(t *testing.T) {
	suite.RunSuite(t, new(AdminServiceSuite))
}

type AdminServiceSuite struct {
	suite.Suite
}

func (s *AdminServiceSuite) TestAdminService_GetAllUsers(t provider.T) {
	userCreator := testobj.NewUserMother()

	t.WithNewStep("success return 2 users", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		users2 := []*models.User{
			userCreator.DefaultUserP(uuid.New()),
			userCreator.DefaultUserP(uuid.New()),
		}

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAll", ctx).Return(users2, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resUsers, err := adminServ.GetAllUsers(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(users2), len(resUsers))
		for i := range users2 {
			sCtx.Assert().Equal(users2[i], resUsers[i])
		}
	}, allure.NewParameter("scenario", "success with 2 users"))

	t.WithNewStep("success return 0 users", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		usersEmpty := make([]*models.User, 0)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAll", ctx).Return(usersEmpty, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resUsers, err := adminServ.GetAllUsers(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Empty(resUsers)
	}, allure.NewParameter("scenario", "success with empty result"))

	t.WithNewStep("error getAll in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		usersEmpty := make([]*models.User, 0)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("userRep error")
		mockUserRep.On("GetAll", ctx).Return(usersEmpty, expectedErr)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllUsers(ctx)

		sCtx.Assert().ErrorIs(err, expectedErr)
	}, allure.NewParameter("scenario", "repository error"))

	t.WithNewStep("error no admin in context", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)

		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("admin error")
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), expectedErr)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllUsers(ctx)

		sCtx.Assert().ErrorIs(err, expectedErr)
	}, allure.NewParameter("scenario", "admin context error"))
}

func (s *AdminServiceSuite) TestAdminService_GetAllEmployees(t provider.T) {
	employeeCreator := testobj.NewEmployeeMother()

	t.WithNewStep("success return 2 employees", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		employees2 := []*models.Employee{
			employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
			employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
		}

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockEmployeeRep.On("GetAll", ctx).Return(employees2, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resEmployees, err := adminServ.GetAllEmployees(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(employees2), len(resEmployees))
		for i := range employees2 {
			sCtx.Assert().Equal(employees2[i], resEmployees[i])
		}
	}, allure.NewParameter("scenario", "success with 2 employees"))

	t.WithNewStep("success return 0 employees", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		employeesEmpty := make([]*models.Employee, 0)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep.On("GetAll", ctx).Return(employeesEmpty, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resUsers, err := adminServ.GetAllEmployees(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Empty(resUsers)
	}, allure.NewParameter("scenario", "success with empty result"))

	t.WithNewStep("error getAll in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		employeesEmpty := make([]*models.Employee, 0)

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		expectedErr := errors.New("emplRep error")
		mockEmployeeRep.On("GetAll", ctx).Return(employeesEmpty, expectedErr)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllEmployees(ctx)

		sCtx.Assert().ErrorIs(err, expectedErr)
	}, allure.NewParameter("scenario", "repository error"))

	t.WithNewStep("error no admin in context", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)

		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("admin error")
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), expectedErr)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllEmployees(ctx)

		sCtx.Assert().ErrorIs(err, expectedErr)
	}, allure.NewParameter("scenario", "admin context error"))
}

func (s *AdminServiceSuite) TestAdminService_ChangeEmployeeRights(t provider.T) {
	employeeCreator := testobj.NewEmployeeMother()

	t.WithNewStep("success", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		valid := true
		employee := employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockEmployeeRep.On("Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)")).Return(employee, nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		err := adminServ.ChangeEmployeeRights(ctx, employee.GetID(), valid)

		sCtx.Assert().NoError(err)
		mockAuthZRep.AssertCalled(t, "AdminIDFromContext", ctx)
		mockEmployeeRep.AssertCalled(t, "Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)"))
	}, allure.NewParameter("scenario", "successful change"))

	t.WithNewStep("error no admin in context", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("admin error")
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), expectedErr)

		valid := true
		employee := employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		err := adminServ.ChangeEmployeeRights(ctx, employee.GetID(), valid)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockAuthZRep.AssertCalled(t, "AdminIDFromContext", ctx)
	}, allure.NewParameter("scenario", "admin context error"))

	t.WithNewStep("error update", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil)

		valid := true
		employee := employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		expectedErr := errors.New("error")
		mockEmployeeRep.On("Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)")).Return(employee, expectedErr)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		err := adminServ.ChangeEmployeeRights(ctx, employee.GetID(), valid)

		sCtx.Assert().ErrorIs(err, expectedErr)
		mockAuthZRep.AssertCalled(t, "AdminIDFromContext", ctx)
		mockEmployeeRep.AssertCalled(t, "Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)"))
	}, allure.NewParameter("scenario", "update error"))
}
