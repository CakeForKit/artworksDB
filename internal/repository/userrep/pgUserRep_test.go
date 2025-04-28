package userrep_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
	"github.com/stateio/testify/assert"
	"github.com/stateio/testify/require"
)

func createTestUser(num int, id uuid.UUID, subscribed bool) *models.User {
	user, _ := models.NewUser(
		id,
		fmt.Sprintf("testUser%d", num),
		fmt.Sprintf("testLogin%d", num),
		fmt.Sprintf("testHashedPassword%d", num),
		time.Now(),
		fmt.Sprintf("user%d@test.com", num),
		subscribed,
	)
	return &user
}

func createTestUsersArr(cnt int) []*models.User {
	users := make([]*models.User, 0)
	for i := 0; i < cnt; i++ {
		users = append(users, createTestUser(i, uuid.New(), true))
	}
	return users
}

func compareUsersWithoutOrder(users1, users2 []*models.User) bool {
	if len(users1) != len(users2) {
		return false
	}
	// Создаем map для учета совпадений
	matched := make([]bool, len(users2))

	for _, u1 := range users1 {
		found := false
		for j, u2 := range users2 {
			if !matched[j] && reflect.DeepEqual(u1, u2) {
				matched[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestUserRep_CRUD(t *testing.T) {
	ctx := context.Background()
	fmt.Printf("HERE 0\n")

	dbCnfg := cnfg.GetTestDatebaseConfig()

	_, pgCreds, err := pgtest.GetTestPostgres(ctx)
	require.NoError(t, err)
	// // defer container.Terminate(ctx)

	urep, err := userrep.NewPgUserRep(ctx, &pgCreds, dbCnfg)
	require.NoError(t, err)

	tests := []struct {
		name          string
		gotUsers      []*models.User
		expectedError error
	}{
		{
			name:          "3 users",
			gotUsers:      createTestUsersArr(3),
			expectedError: nil,
		},
		{
			name:          "1 user",
			gotUsers:      createTestUsersArr(1),
			expectedError: nil,
		},
		{
			name:          "0 user",
			gotUsers:      make([]*models.User, 0),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = pgtest.MigrateUp(ctx)
			require.NoError(t, err)

			for _, u := range tt.gotUsers {
				err = urep.Add(ctx, u)
				require.NoError(t, err)
			}

			resUsers, err := urep.GetAll(ctx)
			require.NoError(t, err)
			assert.True(t, len(resUsers) == len(tt.gotUsers))
			for i, _ := range resUsers {
				assert.True(t, models.CmpUsers(resUsers[i], tt.gotUsers[i]))
			}

			for _, u := range tt.gotUsers {
				resUser, err := urep.GetByID(ctx, u.GetID())
				require.NoError(t, err)

				assert.True(t, models.CmpUsers(resUser, u))
			}

			for _, u := range tt.gotUsers {
				resUser, err := urep.GetByLogin(ctx, u.GetLogin())
				require.NoError(t, err)

				assert.True(t, models.CmpUsers(resUser, u))
			}

			for _, u := range tt.gotUsers {
				err := urep.Delete(ctx, u.GetID())
				require.NoError(t, err)
			}
			resUsers, err = urep.GetAll(ctx)
			require.NoError(t, err)
			assert.True(t, len(resUsers) == 0)

			err = pgtest.MigrateDown(ctx)
			require.NoError(t, err)
		})
	}

	t.Run("Update user fields", func(t *testing.T) {
		err = pgtest.MigrateUp(ctx)
		require.NoError(t, err)

		testUser := createTestUser(0, uuid.New(), true)
		err = urep.Add(ctx, testUser)
		require.NoError(t, err)

		newUsername := "updatedUsername"
		newEmail := "updated@email.com"

		updateFunc := func(u *models.User) (*models.User, error) {
			newUser, err := models.NewUser(
				u.GetID(),
				newUsername,
				u.GetLogin(),
				u.GetHashedPassword(),
				u.GetCreatedAt(),
				newEmail,
				u.IsSubscribedToMail(),
			)
			if err != nil {
				return nil, err
			}
			return &newUser, nil
		}
		updatedUser, err := urep.Update(ctx, testUser.GetID(), updateFunc)
		require.NoError(t, err)

		// Проверяем обновленные поля
		assert.Equal(t, newUsername, updatedUser.GetUsername())
		assert.Equal(t, newEmail, updatedUser.GetEmail())

		// Проверяем, что остальные поля не изменились
		assert.Equal(t, testUser.GetLogin(), updatedUser.GetLogin())
		assert.Equal(t, testUser.GetHashedPassword(), updatedUser.GetHashedPassword())
		assert.Equal(t, testUser.IsSubscribedToMail(), updatedUser.IsSubscribedToMail())

		// Проверяем через GetByID
		fetchedUser, err := urep.GetByID(ctx, testUser.GetID())
		require.NoError(t, err)
		assert.Equal(t, newUsername, fetchedUser.GetUsername())
		assert.Equal(t, newEmail, fetchedUser.GetEmail())

		err = pgtest.MigrateDown(ctx)
		require.NoError(t, err)
	})

	t.Run("Update subscription status", func(t *testing.T) {
		err = pgtest.MigrateUp(ctx)
		require.NoError(t, err)
		defer func() {
			err = pgtest.MigrateDown(ctx)
			require.NoError(t, err)
		}()

		// Создаем тестового пользователя с подпиской
		testUser := createTestUser(0, uuid.New(), true)
		err = urep.Add(ctx, testUser)
		require.NoError(t, err)

		// Тест 1: Отключаем подписку
		newStatus := false
		err = urep.UpdateSubscribeToMailing(ctx, testUser.GetID(), newStatus)
		require.NoError(t, err)

		// Проверяем обновление
		updatedUser, err := urep.GetByID(ctx, testUser.GetID())
		require.NoError(t, err)
		assert.Equal(t, newStatus, updatedUser.IsSubscribedToMail())

		// Проверяем, что другие поля не изменились
		assert.Equal(t, testUser.GetUsername(), updatedUser.GetUsername())
		assert.Equal(t, testUser.GetEmail(), updatedUser.GetEmail())
		assert.Equal(t, testUser.GetLogin(), updatedUser.GetLogin())

		// Тест 2: Включаем подписку обратно
		err = urep.UpdateSubscribeToMailing(ctx, testUser.GetID(), true)
		require.NoError(t, err)

		// Проверяем обновление
		updatedUser, err = urep.GetByID(ctx, testUser.GetID())
		require.NoError(t, err)
		assert.True(t, updatedUser.IsSubscribedToMail())

		// Тест 3: Попытка обновления несуществующего пользователя
		nonExistentID := uuid.New()
		err = urep.UpdateSubscribeToMailing(ctx, nonExistentID, true)
		assert.Error(t, err)                                    // Ожидаем ошибку
		assert.True(t, errors.Is(err, userrep.ErrRowsAffected)) // Проверяем тип ошибки
	})

}
