package adminrep_test

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
	ctx     context.Context
	arep    *adminrep.PgAdminRep
	dbCnfg  *cnfg.DatebaseConfig
	pgCreds *cnfg.DatebaseCredentials
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {
		dbCnfg := cnfg.GetTestDatebaseConfig()
		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		arep, err := adminrep.NewPgAdminRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:     ctx,
			arep:    arep,
			dbCnfg:  dbCnfg,
			pgCreds: &pgCreds,
		}
	})
	pgTestConfig := cnfg.GetPgTestConfig()

	err := pgtest.MigrateUp(ctx, pgTestConfig.MigrationDir, th.pgCreds)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx, pgTestConfig.MigrationDir, th.pgCreds)
		require.NoError(t, err)
	})

	return th
}

func (th *testHelper) createTestAdmin(num int) *models.Admin {
	admin, err := models.NewAdmin(
		uuid.New(),
		fmt.Sprintf("admin%d", num),
		fmt.Sprintf("adminL%d", num),
		fmt.Sprintf("hashPadmin%d", num),
		time.Now().UTC().Truncate(time.Microsecond),
		true,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestAdmin failed: %v", err))
	}
	return &admin
}

func (th *testHelper) createAndAddAdmin(t *testing.T, num int) *models.Admin {
	admin := th.createTestAdmin(num)
	err := th.arep.Add(th.ctx, admin)
	require.NoError(t, err)
	return admin
}

func TestAdminRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() []*models.Admin
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			setup:         func() []*models.Admin { return nil },
			wantLen:       0,
			expectedError: nil,
		},
		{
			name: "Should return all admins",
			setup: func() []*models.Admin {
				admins := make([]*models.Admin, 3)
				for i := range admins {
					admins[i] = th.createAndAddAdmin(t, i)
				}
				return admins
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedAdmins := tt.setup()

			admins, err := th.arep.GetAll(th.ctx)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Len(t, admins, tt.wantLen)

			if len(expectedAdmins) > 0 {
				for i, expected := range expectedAdmins {
					assert.Equal(t, expected.GetID(), admins[i].GetID())
					assert.Equal(t, expected.GetUsername(), admins[i].GetUsername())
					assert.Equal(t, expected.GetLogin(), admins[i].GetLogin())
				}
			}
		})
	}
}

func TestAdminRep_GetByID(t *testing.T) {
	th := setupTestHelper(t)

	admin := th.createAndAddAdmin(t, 1)
	nonExistentID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *models.Admin
		wantErr error
	}{
		{
			name: "Should return admin by ID",
			id:   admin.GetID(),
			want: admin,
		},
		{
			name:    "Should return error for non-existent ID",
			id:      nonExistentID,
			wantErr: adminrep.ErrAdminNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := th.arep.GetByID(th.ctx, tt.id)

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

func TestAdminRep_GetByLogin(t *testing.T) {
	th := setupTestHelper(t)

	admin := th.createAndAddAdmin(t, 1)

	t.Run("Should return admin by login", func(t *testing.T) {
		got, err := th.arep.GetByLogin(th.ctx, admin.GetLogin())
		require.NoError(t, err)
		assert.Equal(t, admin.GetID(), got.GetID())
		assert.Equal(t, admin.GetUsername(), got.GetUsername())
		assert.Equal(t, admin.GetLogin(), got.GetLogin())
	})

	t.Run("Should return error for non-existent login", func(t *testing.T) {
		_, err := th.arep.GetByLogin(th.ctx, "non_existent_login")
		assert.ErrorIs(t, err, adminrep.ErrAdminNotFound)
	})
}

func TestAdminRep_Update(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		updateFunc func(*models.Admin) (*models.Admin, error)
		wantCheck  func(*testing.T, *models.Admin)
		wantErr    bool
	}{
		{
			name: "Should update admin",
			updateFunc: func(a *models.Admin) (*models.Admin, error) {
				newAdmin, err := models.NewAdmin(
					a.GetID(),
					"new_username",
					"new_login",
					"new_hashed_password",
					a.GetCreatedAt(),
					false,
				)
				return &newAdmin, err
			},
			wantCheck: func(t *testing.T, a *models.Admin) {
				assert.Equal(t, "new_username", a.GetUsername())
				assert.Equal(t, "new_login", a.GetLogin())
				assert.False(t, a.IsValid())
			},
		},
		{
			name: "Should return error from updateFunc",
			updateFunc: func(a *models.Admin) (*models.Admin, error) {
				return nil, errors.New("update error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			admin := th.createAndAddAdmin(t, 1)

			err := th.arep.Update(th.ctx, admin.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			dbAdmin, err := th.arep.GetByID(th.ctx, admin.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbAdmin)
		})
	}
}
