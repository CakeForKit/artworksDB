package employeerep_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	th     *testHelper
	pgOnce sync.Once
)

type testHelper struct {
	ctx          context.Context
	erep         *employeerep.PgEmployeeRep
	dbCnfg       *cnfg.DatebaseConfig
	pgTestConfig *cnfg.PostgresTestConfig
	pgCreds      *cnfg.PostgresCredentials
	adminID      uuid.UUID
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {
		dbCnfg := cnfg.GetTestDatebaseConfig()
		pgTestConfig := cnfg.GetPgTestConfig()

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		erep, err := employeerep.NewPgEmployeeRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:          ctx,
			erep:         erep,
			adminID:      uuid.New(),
			dbCnfg:       dbCnfg,
			pgTestConfig: pgTestConfig,
			pgCreds:      &pgCreds,
		}
	})
	err := pgtest.MigrateUp(ctx, th.pgTestConfig, th.pgCreds)
	require.NoError(t, err)

	admin, err := models.NewAdmin(
		th.adminID,
		"admin",
		"adminL",
		"hashPadmin",
		time.Now().UTC().Truncate(time.Microsecond),
		true,
	)
	require.NoError(t, err)
	arep, err := adminrep.NewPgAdminRep(ctx, th.pgCreds, th.dbCnfg)
	require.NoError(t, err)
	err = arep.Add(ctx, &admin)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx, th.pgTestConfig, th.pgCreds)
		require.NoError(t, err)
	})

	return th

}

func (th *testHelper) createTestEmployee(num int) *models.Employee {
	employee, err := models.NewEmployee(
		uuid.New(),
		fmt.Sprintf("testEmployee%d", num),
		fmt.Sprintf("testLogin%d", num),
		fmt.Sprintf("testHashedPassword%d", num),
		time.Now().UTC().Truncate(time.Microsecond),
		true,
		th.adminID,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestEmployee failed: %v", err))
	}
	return &employee
}

func (th *testHelper) createAndAddEmployee(t *testing.T, num int) *models.Employee {
	employee := th.createTestEmployee(num)
	// fmt.Printf("EMPLOYEE: %+v\n", employee)
	err := th.erep.Add(th.ctx, employee)
	require.NoError(t, err)
	return employee
}

func TestEmployeeRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() []*models.Employee
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			setup:         func() []*models.Employee { return nil },
			wantLen:       0,
			expectedError: employeerep.ErrEmployeeNotFound,
		},
		{
			name: "Should return all employees",
			setup: func() []*models.Employee {
				employees := make([]*models.Employee, 3)
				for i := range employees {
					employees[i] = th.createAndAddEmployee(t, i)
				}
				return employees
			},
			wantLen:       3,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedEmployees := tt.setup()

			employees, err := th.erep.GetAll(th.ctx)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Len(t, employees, tt.wantLen)

			if len(expectedEmployees) > 0 {
				for i, expected := range expectedEmployees {
					assert.Equal(t, expected.GetID(), employees[i].GetID())
					assert.Equal(t, expected.GetUsername(), employees[i].GetUsername())
					assert.Equal(t, expected.GetLogin(), employees[i].GetLogin())
				}
			}
		})
	}
}

func TestEmployeeRep_GetByID(t *testing.T) {
	th := setupTestHelper(t)

	employee := th.createAndAddEmployee(t, 1)
	nonExistentID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *models.Employee
		wantErr error
	}{
		{
			name: "Should return employee by ID",
			id:   employee.GetID(),
			want: employee,
		},
		{
			name:    "Should return error for non-existent ID",
			id:      nonExistentID,
			wantErr: employeerep.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := th.erep.GetByID(th.ctx, tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.GetID(), got.GetID())
			assert.Equal(t, tt.want.GetUsername(), got.GetUsername())
			assert.Equal(t, tt.want.GetLogin(), got.GetLogin())
		})
	}
}

func TestEmployeeRep_GetByLogin(t *testing.T) {
	th := setupTestHelper(t)

	employee := th.createAndAddEmployee(t, 1)

	t.Run("Should return employee by login", func(t *testing.T) {
		got, err := th.erep.GetByLogin(th.ctx, employee.GetLogin())
		require.NoError(t, err)
		assert.Equal(t, employee.GetID(), got.GetID())
		assert.Equal(t, employee.GetUsername(), got.GetUsername())
		assert.Equal(t, employee.GetLogin(), got.GetLogin())
	})

	t.Run("Should return error for non-existent login", func(t *testing.T) {
		res, err := th.erep.GetByLogin(th.ctx, "non_existent_login")
		assert.ErrorIs(t, err, employeerep.ErrEmployeeNotFound)
		assert.Nil(t, res)
	})
}

func TestEmployeeRep_Add(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should add new employee", func(t *testing.T) {
		employee := th.createTestEmployee(1)
		err := th.erep.Add(th.ctx, employee)
		require.NoError(t, err)

		got, err := th.erep.GetByID(th.ctx, employee.GetID())
		require.NoError(t, err)
		assert.Equal(t, employee.GetID(), got.GetID())
		assert.Equal(t, employee.GetUsername(), got.GetUsername())
		assert.Equal(t, employee.GetLogin(), got.GetLogin())
	})

	t.Run("Should return error for duplicate login", func(t *testing.T) {
		employee1 := th.createAndAddEmployee(t, 3)
		num := 2
		employee2, err := models.NewEmployee(
			uuid.New(),
			fmt.Sprintf("testEmployee%d", num),
			employee1.GetLogin(),
			fmt.Sprintf("testHashedPassword%d", num),
			time.Now().UTC().Truncate(time.Microsecond),
			true,
			th.adminID,
		)
		assert.NoError(t, err)
		// employee2 := th.createTestEmployee(2)
		// employee2.SetLogin(employee1.GetLogin())

		err = th.erep.Add(th.ctx, &employee2)
		assert.Error(t, err)
	})
}

func TestEmployeeRep_Delete(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should delete existing employee", func(t *testing.T) {
		employee := th.createAndAddEmployee(t, 1)

		err := th.erep.Delete(th.ctx, employee.GetID())
		require.NoError(t, err)

		_, err = th.erep.GetByID(th.ctx, employee.GetID())
		assert.ErrorIs(t, err, employeerep.ErrEmployeeNotFound)
	})

	t.Run("Should return error for non-existent employee", func(t *testing.T) {
		err := th.erep.Delete(th.ctx, uuid.New())
		assert.ErrorIs(t, err, employeerep.ErrRowsAffected)
	})
}

func TestEmployeeRep_Update(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		updateFunc func(*models.Employee) (*models.Employee, error)
		wantCheck  func(*testing.T, *models.Employee)
		wantErr    bool
	}{
		{
			name: "Should update username",
			updateFunc: func(e *models.Employee) (*models.Employee, error) {
				newEmployee, err := models.NewEmployee(
					e.GetID(),
					"new_username",
					"new_login",
					e.GetHashedPassword(),
					e.GetCreatedAt(),
					true,
					th.adminID,
				)
				return &newEmployee, err
			},
			wantCheck: func(t *testing.T, e *models.Employee) {
				assert.Equal(t, "new_username", e.GetUsername())
				assert.Equal(t, "new_login", e.GetLogin())
			},
		},
		{
			name: "Should return error from updateFunc",
			updateFunc: func(e *models.Employee) (*models.Employee, error) {
				return nil, errors.New("update error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employee := th.createAndAddEmployee(t, 1)

			updated, err := th.erep.Update(th.ctx, employee.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.wantCheck(t, updated)

			dbEmployee, err := th.erep.GetByID(th.ctx, employee.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbEmployee)
		})
	}
}
