package userrep_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
	"github.com/stateio/testify/require"
	"github.com/stretchr/testify/assert"
)

// testHelper содержит общие методы для тестов
type testHelper struct {
	ctx    context.Context
	urep   *userrep.PgUserRep
	dbCnfg *cnfg.DatebaseConfig
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	dbCnfg := cnfg.GetTestDatebaseConfig()

	_, pgCreds, err := pgtest.GetTestPostgres(ctx)
	require.NoError(t, err)

	urep, err := userrep.NewPgUserRep(ctx, &pgCreds, dbCnfg)
	require.NoError(t, err)

	err = pgtest.MigrateUp(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx)
		require.NoError(t, err)
	})

	return &testHelper{
		ctx:    ctx,
		urep:   urep,
		dbCnfg: dbCnfg,
	}
}

func (th *testHelper) createTestUser(num int, subscribed bool) *models.User {
	user, err := models.NewUser(
		uuid.New(),
		fmt.Sprintf("testUser%d", num),
		fmt.Sprintf("testLogin%d", num),
		fmt.Sprintf("testHashedPassword%d", num),
		time.Now().UTC().Truncate(time.Microsecond), // Нормализация времени
		fmt.Sprintf("user%d@test.com", num),
		subscribed,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestUser failed: %v", err))
	}
	return &user
}

func (th *testHelper) createAndAddUser(t *testing.T, num int, subscribed bool) *models.User {
	user := th.createTestUser(num, subscribed)
	err := th.urep.Add(th.ctx, user)
	require.NoError(t, err)
	return user
}

func TestUserRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() []*models.User
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			setup:         func() []*models.User { return nil },
			wantLen:       0,
			expectedError: userrep.ErrUserNotFound,
		},
		{
			name: "Should return all users",
			setup: func() []*models.User {
				users := make([]*models.User, 3)
				for i := range users {
					users[i] = th.createAndAddUser(t, i, true)
				}
				return users
			},
			wantLen:       3,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedUsers := tt.setup()

			users, err := th.urep.GetAll(th.ctx)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Len(t, users, tt.wantLen)

			if len(expectedUsers) > 0 {
				for i, expected := range expectedUsers {
					assert.True(t, models.CmpUsers(expected, users[i]),
						"User %d mismatch", i)
				}
			}
		})
	}
}

func TestUserRep_GetAllSubscribed(t *testing.T) {
	th := setupTestHelper(t)

	// Создаем тестовых пользователей
	subscribedUsers := []*models.User{
		th.createAndAddUser(t, 1, true),
		th.createAndAddUser(t, 2, true),
	}
	th.createAndAddUser(t, 3, false) // Не подписан

	t.Run("Should return only subscribed users", func(t *testing.T) {
		users, err := th.urep.GetAllSubscribed(th.ctx)
		require.NoError(t, err)

		assert.Len(t, users, len(subscribedUsers))
		for i, user := range subscribedUsers {
			assert.True(t, models.CmpUsers(user, users[i]),
				"Subscribed user %d mismatch", i)
		}
	})
}

func TestUserRep_GetByID(t *testing.T) {
	th := setupTestHelper(t)

	user := th.createAndAddUser(t, 1, true)
	nonExistentID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *models.User
		wantErr error
	}{
		{
			name: "Should return user by ID",
			id:   user.GetID(),
			want: user,
		},
		{
			name:    "Should return error for non-existent ID",
			id:      nonExistentID,
			wantErr: userrep.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := th.urep.GetByID(th.ctx, tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.True(t, models.CmpUsers(tt.want, got))
		})
	}
}

func TestUserRep_GetByLogin(t *testing.T) {
	th := setupTestHelper(t)

	user := th.createAndAddUser(t, 1, true)

	t.Run("Should return user by login", func(t *testing.T) {
		got, err := th.urep.GetByLogin(th.ctx, user.GetLogin())
		require.NoError(t, err)
		assert.True(t, models.CmpUsers(user, got))
	})

	t.Run("Should return error for non-existent login", func(t *testing.T) {
		res, err := th.urep.GetByLogin(th.ctx, "non_existent_login")
		assert.ErrorIs(t, err, userrep.ErrUserNotFound)
		assert.Nil(t, res)
	})
}

func TestUserRep_Add(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should add new user", func(t *testing.T) {
		user := th.createTestUser(1, true)
		err := th.urep.Add(th.ctx, user)
		require.NoError(t, err)

		// Проверяем, что пользователь действительно добавлен
		got, err := th.urep.GetByID(th.ctx, user.GetID())
		require.NoError(t, err)
		assert.True(t, models.CmpUsers(user, got))
	})

	t.Run("Should return error for duplicate login", func(t *testing.T) {
		user1 := th.createAndAddUser(t, 3, true)
		// user2 := th.createTestUser(2, true)
		num := 2
		user2, err := models.NewUser(
			uuid.New(),
			fmt.Sprintf("testUser%d", num),
			user1.GetLogin(),
			fmt.Sprintf("testHashedPassword%d", num),
			time.Now().UTC().Truncate(time.Microsecond),
			fmt.Sprintf("user%d@test.com", num),
			true,
		)
		assert.NoError(t, err)

		err = th.urep.Add(th.ctx, &user2)
		assert.Error(t, err)
	})
}

func TestUserRep_Delete(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should delete existing user", func(t *testing.T) {
		user := th.createAndAddUser(t, 1, true)

		err := th.urep.Delete(th.ctx, user.GetID())
		require.NoError(t, err)

		// Проверяем, что пользователь удален
		_, err = th.urep.GetByID(th.ctx, user.GetID())
		assert.ErrorIs(t, err, userrep.ErrUserNotFound)
	})

	t.Run("Should return error for non-existent user", func(t *testing.T) {
		err := th.urep.Delete(th.ctx, uuid.New())
		assert.ErrorIs(t, err, userrep.ErrRowsAffected)
	})
}

func TestUserRep_Update(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		updateFunc func(*models.User) (*models.User, error)
		wantCheck  func(*testing.T, *models.User)
		wantErr    bool
	}{
		{
			name: "Should update username",
			updateFunc: func(u *models.User) (*models.User, error) {
				newUser, err := models.NewUser(
					u.GetID(),
					"new_username",
					"new_login",
					u.GetHashedPassword(),
					u.GetCreatedAt(),
					"new_email@test.ru",
					u.IsSubscribedToMail(),
				)
				return &newUser, err
			},
			wantCheck: func(t *testing.T, u *models.User) {
				assert.Equal(t, "new_username", u.GetUsername())
				assert.Equal(t, "new_login", u.GetLogin())
				assert.Equal(t, "new_email@test.ru", u.GetEmail())
			},
		},
		{
			name: "Should return error from updateFunc",
			updateFunc: func(u *models.User) (*models.User, error) {
				return nil, errors.New("update error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := th.createAndAddUser(t, 1, true)

			updated, err := th.urep.Update(th.ctx, user.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.wantCheck(t, updated)

			// Проверяем, что изменения сохранились в БД
			dbUser, err := th.urep.GetByID(th.ctx, user.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbUser)
		})
	}
}

func TestUserRep_UpdateSubscribeToMailing(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		initialSub bool
		newSub     bool
		wantSub    bool
		wantErr    bool
	}{
		{"Enable subscription", false, true, true, false},
		{"Disable subscription", true, false, false, false},
		{"No change needed", true, true, true, false},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := th.createAndAddUser(t, i, tt.initialSub)

			err := th.urep.UpdateSubscribeToMailing(th.ctx, user.GetID(), tt.newSub)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Проверяем обновленное состояние
			updated, err := th.urep.GetByID(th.ctx, user.GetID())
			require.NoError(t, err)
			assert.Equal(t, tt.wantSub, updated.IsSubscribedToMail())

			// Проверяем, что другие поля не изменились
			assert.Equal(t, user.GetUsername(), updated.GetUsername())
			assert.Equal(t, user.GetEmail(), updated.GetEmail())
		})
	}

	t.Run("Should return error for non-existent user", func(t *testing.T) {
		err := th.urep.UpdateSubscribeToMailing(th.ctx, uuid.New(), true)
		assert.ErrorIs(t, err, userrep.ErrRowsAffected)
	})
}
