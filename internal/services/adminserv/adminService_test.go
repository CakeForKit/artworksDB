package adminserv_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/adminserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAdminService_GetAllUsers(t *testing.T) {
	userCreator := testobj.NewUserMother()

	t.Run("success return 2 users", func(t *testing.T) {
		ctx := context.Background()

		users2 := []*models.User{
			userCreator.DefaultUserP(uuid.New()),
			userCreator.DefaultUserP(uuid.New()),
		}

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAll", ctx).Return(users2, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resUsers, err := adminServ.GetAllUsers(ctx)
		require.Nil(t, err)
		require.True(t, len(users2) == len(resUsers))
		for i := range len(users2) {
			require.True(t, users2[i] == resUsers[i])
		}
	})
	t.Run("success return 0 users", func(t *testing.T) {
		ctx := context.Background()

		usersEmpty := make([]*models.User, 0)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		mockUserRep.On("GetAll", ctx).Return(usersEmpty, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resUsers, err := adminServ.GetAllUsers(ctx)
		require.Nil(t, err)
		require.True(t, len(resUsers) == 0)
	})

	t.Run("error getAll in rep", func(t *testing.T) {
		ctx := context.Background()

		usersEmpty := make([]*models.User, 0)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		expectedErr := errors.New("userRep error")
		mockUserRep.On("GetAll", ctx).Return(usersEmpty, expectedErr)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllUsers(ctx)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("error no admin in context", func(t *testing.T) {
		ctx := context.Background()

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)

		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("admin error")
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), expectedErr) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllUsers(ctx)

		require.ErrorIs(t, err, expectedErr)
	})
}

func TestAdminService_GetAllEmployees(t *testing.T) {
	employeeCreator := testobj.NewEmployeeMother()

	t.Run("success return 2 employees", func(t *testing.T) {
		ctx := context.Background()

		employees2 := []*models.Employee{
			employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
			employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New()),
		}

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockEmployeeRep.On("GetAll", ctx).Return(employees2, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resEmployees, err := adminServ.GetAllEmployees(ctx)
		require.Nil(t, err)
		require.True(t, len(employees2) == len(resEmployees))
		for i := range len(employees2) {
			require.True(t, employees2[i] == resEmployees[i])
		}
	})

	t.Run("success return 0 employees", func(t *testing.T) {
		ctx := context.Background()

		employeesEmpty := make([]*models.Employee, 0)

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep.On("GetAll", ctx).Return(employeesEmpty, nil)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		resUsers, err := adminServ.GetAllEmployees(ctx)
		require.Nil(t, err)
		require.True(t, len(resUsers) == 0)
	})

	t.Run("error getAll in rep", func(t *testing.T) {
		ctx := context.Background()

		employeesEmpty := make([]*models.Employee, 0)

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		expectedErr := errors.New("emplRep error")
		mockEmployeeRep.On("GetAll", ctx).Return(employeesEmpty, expectedErr)

		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllEmployees(ctx)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("error no admin in context", func(t *testing.T) {
		ctx := context.Background()

		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockUserRep := new(userrep.MockUserRep)

		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("admin error")
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), expectedErr) // adminID не важен

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		_, err := adminServ.GetAllEmployees(ctx)

		require.ErrorIs(t, err, expectedErr)
	})

}

func TestAdminService_ChangeEmployeeRights(t *testing.T) {
	employeeCreator := testobj.NewEmployeeMother()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		valid := true
		employee := employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		mockEmployeeRep.On("Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)")).Return(employee, nil)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		err := adminServ.ChangeEmployeeRights(ctx, employee.GetID(), valid)
		require.Nil(t, err)
		mockAuthZRep.AssertCalled(t, "AdminIDFromContext", ctx)
		mockEmployeeRep.AssertCalled(t, "Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)"))
	})
	t.Run("error no admin in context", func(t *testing.T) {
		ctx := context.Background()
		mockAuthZRep := new(auth.MockAuthZ)
		expectedErr := errors.New("admin error")
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), expectedErr) // adminID не важен

		valid := true
		employee := employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		err := adminServ.ChangeEmployeeRights(ctx, employee.GetID(), valid)
		require.ErrorIs(t, err, expectedErr)
		mockAuthZRep.AssertCalled(t, "AdminIDFromContext", ctx)
	})

	t.Run("error update", func(t *testing.T) {
		ctx := context.Background()
		mockAuthZRep := new(auth.MockAuthZ)
		mockAuthZRep.On("AdminIDFromContext", ctx).Return(uuid.New(), nil) // adminID не важен

		valid := true
		employee := employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())

		mockUserRep := new(userrep.MockUserRep)
		mockEmployeeRep := new(employeerep.MockEmployeeRep)
		expectedErr := errors.New("error")
		mockEmployeeRep.On("Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)")).Return(employee, expectedErr)

		adminServ := adminserv.NewAdminService(mockEmployeeRep, mockUserRep, mockAuthZRep)
		// ACT
		err := adminServ.ChangeEmployeeRights(ctx, employee.GetID(), valid)
		require.ErrorIs(t, err, expectedErr)
		mockAuthZRep.AssertCalled(t, "AdminIDFromContext", ctx)
		mockEmployeeRep.AssertCalled(t, "Update", ctx, employee.GetID(), mock.AnythingOfType("func(*models.Employee) (*models.Employee, error)"))
	})
}
