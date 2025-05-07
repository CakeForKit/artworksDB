package eventrep_test

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
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	th     *testHelper
	pgOnce sync.Once
)

type testHelper struct {
	ctx          context.Context
	erep         *eventrep.PgEventRep
	dbCnfg       *cnfg.DatebaseConfig
	pgTestConfig *cnfg.PostgresTestConfig
	pgCreds      *cnfg.PostgresCredentials
	employeeID   uuid.UUID
}

func addEmployee(t *testing.T, ctx context.Context, employeeID uuid.UUID, pgCreds *cnfg.PostgresCredentials, dbCnfg *cnfg.DatebaseConfig) {
	admin, err := models.NewAdmin(
		uuid.New(),
		"admin",
		"adminL",
		"hashPadmin",
		time.Now().UTC().Truncate(time.Microsecond),
		true,
	)
	// fmt.Printf("ADMIN: %+v\n", admin)
	require.NoError(t, err)
	arep, err := adminrep.NewPgAdminRep(ctx, pgCreds, dbCnfg)
	require.NoError(t, err)
	err = arep.Add(ctx, &admin)
	require.NoError(t, err)

	employee, err := models.NewEmployee(
		employeeID,
		"empTest",
		"loginTest",
		"hpTest",
		time.Now().UTC().Truncate(time.Microsecond),
		true,
		admin.GetID(),
	)
	require.NoError(t, err)
	erep, err := employeerep.NewPgEmployeeRep(ctx, pgCreds, dbCnfg)
	require.NoError(t, err)
	err = erep.Add(ctx, &employee)
	require.NoError(t, err)
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {
		dbCnfg := cnfg.GetTestDatebaseConfig()
		pgTestConfig := cnfg.GetPgTestConfig()

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		erep, err := eventrep.NewPgEventRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:          ctx,
			erep:         erep,
			employeeID:   uuid.New(),
			dbCnfg:       dbCnfg,
			pgTestConfig: pgTestConfig,
			pgCreds:      &pgCreds,
		}
	})
	err := pgtest.MigrateUp(ctx, th.pgTestConfig, th.pgCreds)
	require.NoError(t, err)
	addEmployee(t, ctx, th.employeeID, th.pgCreds, th.dbCnfg)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx, th.pgTestConfig, th.pgCreds)
		require.NoError(t, err)
	})

	return th
}

func (th *testHelper) createTestEvent(num int) *models.Event {
	event, err := models.NewEvent(
		uuid.New(),
		fmt.Sprintf("Event %d", num),
		time.Now().AddDate(0, 0, num),
		time.Now().AddDate(0, 0, num+1),
		fmt.Sprintf("Address %d", num),
		true,
		th.employeeID,
		100,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestEvent failed: %v", err))
	}
	return &event
}

func (th *testHelper) createAndAddEvent(t *testing.T, num int) *models.Event {
	event := th.createTestEvent(num)
	err := th.erep.Add(th.ctx, event)
	require.NoError(t, err)
	return event
}

func TestEventRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() []*models.Event
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			setup:         func() []*models.Event { return nil },
			wantLen:       0,
			expectedError: eventrep.ErrEventNotFound,
		},
		{
			name: "Should return all events",
			setup: func() []*models.Event {
				events := make([]*models.Event, 3)
				for i := range events {
					events[i] = th.createAndAddEvent(t, i)
				}
				return events
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedEvents := tt.setup()

			events, err := th.erep.GetAll(th.ctx)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Len(t, events, tt.wantLen)

			if len(expectedEvents) > 0 {
				for i, expected := range expectedEvents {
					assert.Equal(t, expected.GetID(), events[i].GetID())
					assert.Equal(t, expected.GetTitle(), events[i].GetTitle())
				}
			}
		})
	}
}

func TestEventRep_GetByID(t *testing.T) {
	th := setupTestHelper(t)

	event := th.createAndAddEvent(t, 1)
	nonExistentID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *models.Event
		wantErr error
	}{
		{
			name: "Should return event by ID",
			id:   event.GetID(),
			want: event,
		},
		{
			name:    "Should return error for non-existent ID",
			id:      nonExistentID,
			wantErr: eventrep.ErrEventNotFound,
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
			assert.Equal(t, tt.want.GetTitle(), got.GetTitle())
		})
	}
}

func TestEventRep_GetByDate(t *testing.T) {
	th := setupTestHelper(t)

	now := time.Now()
	event := th.createAndAddEvent(t, 1)

	tests := []struct {
		name    string
		dateBeg time.Time
		dateEnd time.Time
		wantLen int
		wantErr bool
		wantID  uuid.UUID
	}{
		{
			name:    "Should return events in date range",
			dateBeg: now.AddDate(0, 0, -1),
			dateEnd: now.AddDate(0, 0, 2),
			wantLen: 1,
			wantID:  event.GetID(),
		},
		{
			name:    "Should return empty for out of range",
			dateBeg: now.AddDate(0, 0, 10),
			dateEnd: now.AddDate(0, 0, 20),
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := th.erep.GetByDate(th.ctx, tt.dateBeg, tt.dateEnd)

			if tt.wantLen == 0 {
				assert.ErrorIs(t, err, eventrep.ErrEventNotFound)
				return
			}

			require.NoError(t, err)
			assert.Len(t, events, tt.wantLen)
			if tt.wantLen > 0 {
				assert.Equal(t, tt.wantID, events[0].GetID())
			}
		})
	}
}

func TestEventRep_Add(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should add new event", func(t *testing.T) {
		event := th.createTestEvent(1)

		err := th.erep.Add(th.ctx, event)
		require.NoError(t, err)

		// Проверяем, что event действительно добавлен
		got, err := th.erep.GetByID(th.ctx, event.GetID())
		require.NoError(t, err)
		assert.Equal(t, event.GetID(), got.GetID())
	})

	t.Run("Should return error for duplicate event", func(t *testing.T) {
		event := th.createAndAddEvent(t, 2)

		// Пытаемся добавить тот же event снова
		err := th.erep.Add(th.ctx, event)
		assert.Error(t, err)
	})
}

func TestEventRep_Delete(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should delete existing event", func(t *testing.T) {
		event := th.createAndAddEvent(t, 1)

		err := th.erep.Delete(th.ctx, event.GetID())
		require.NoError(t, err)

		// Проверяем, что event удален
		_, err = th.erep.GetByID(th.ctx, event.GetID())
		assert.ErrorIs(t, err, eventrep.ErrEventNotFound)
	})

	t.Run("Should return error for non-existent event", func(t *testing.T) {
		err := th.erep.Delete(th.ctx, uuid.New())
		assert.ErrorIs(t, err, eventrep.ErrRowsAffected)
	})
}

func TestEventRep_Update(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		updateFunc func(*models.Event) (*models.Event, error)
		wantCheck  func(*testing.T, *models.Event)
		wantErr    bool
	}{
		{
			name: "Should update event",
			updateFunc: func(e *models.Event) (*models.Event, error) {
				newEvent, err := models.NewEvent(
					e.GetID(),
					"Updated Title",
					e.GetDateBegin().AddDate(0, 0, 1),
					e.GetDateEnd().AddDate(0, 0, 1),
					"Updated Address",
					false,
					th.employeeID,
					200,
				)
				return &newEvent, err
			},
			wantCheck: func(t *testing.T, e *models.Event) {
				assert.Equal(t, "Updated Title", e.GetTitle())
				assert.Equal(t, "Updated Address", e.GetAddress())
				assert.Equal(t, 200, e.GetTicketCount())
			},
		},
		{
			name: "Should return error from updateFunc",
			updateFunc: func(e *models.Event) (*models.Event, error) {
				return nil, errors.New("update error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := th.createAndAddEvent(t, 1)

			updated, err := th.erep.Update(th.ctx, event.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.wantCheck(t, updated)

			// Проверяем, что изменения сохранились в БД
			dbEvent, err := th.erep.GetByID(th.ctx, event.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbEvent)
		})
	}
}
