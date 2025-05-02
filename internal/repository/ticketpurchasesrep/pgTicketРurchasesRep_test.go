package ticketpurchasesrep_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testHelper содержит общие методы для тестов
type testHelper struct {
	ctx      context.Context
	tprep    *ticketpurchasesrep.PgTicketPurchasesRep
	dbCnfg   *cnfg.DatebaseConfig
	userIDs  []uuid.UUID
	eventIDs []uuid.UUID
}

func addUser(t *testing.T, ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbCnfg *cnfg.DatebaseConfig, num int) uuid.UUID {
	userID := uuid.New()
	user, err := models.NewUser(
		userID,
		fmt.Sprintf("user%d", num),
		fmt.Sprintf("userL%d", num),
		"hashPuser",
		time.Now().UTC().Truncate(time.Microsecond),
		fmt.Sprintf("user%d@test.ru", num),
		true,
	)
	require.NoError(t, err)
	urep, err := userrep.NewPgUserRep(ctx, pgCreds, dbCnfg)
	require.NoError(t, err)
	err = urep.Add(ctx, &user)
	require.NoError(t, err)

	return userID
}

func addEmployee(t *testing.T, ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbCnfg *cnfg.DatebaseConfig) uuid.UUID {
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
		uuid.New(),
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

	return employee.GetID()
}

func addEvent(t *testing.T, ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbCnfg *cnfg.DatebaseConfig, employeeID uuid.UUID, num int) uuid.UUID {
	event, err := models.NewEvent(
		uuid.New(),
		fmt.Sprintf("Event %d", num),
		time.Now().Add(time.Hour),
		time.Now().Add(2*time.Hour),
		fmt.Sprintf("Address %d", num),
		true,
		employeeID,
		100+num,
	)
	require.NoError(t, err)
	erep, err := eventrep.NewPgEventRep(ctx, pgCreds, dbCnfg)
	require.NoError(t, err)
	err = erep.Add(ctx, &event)
	require.NoError(t, err)

	return event.GetID()
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	dbCnfg := cnfg.GetTestDatebaseConfig()

	_, pgCreds, err := pgtest.GetTestPostgres(ctx)
	require.NoError(t, err)

	err = pgtest.MigrateUp(ctx)
	require.NoError(t, err)

	employeeID := addEmployee(t, ctx, &pgCreds, dbCnfg)
	userID1 := addUser(t, ctx, &pgCreds, dbCnfg, 1)
	userID2 := addUser(t, ctx, &pgCreds, dbCnfg, 2)
	userID3 := addUser(t, ctx, &pgCreds, dbCnfg, 3)
	eventID1 := addEvent(t, ctx, &pgCreds, dbCnfg, employeeID, 1)
	eventID2 := addEvent(t, ctx, &pgCreds, dbCnfg, employeeID, 2)
	eventID3 := addEvent(t, ctx, &pgCreds, dbCnfg, employeeID, 3)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx)
		require.NoError(t, err)
	})

	tprep, err := ticketpurchasesrep.NewPgTicketPurchasesRep(ctx, &pgCreds, dbCnfg)
	require.NoError(t, err)

	return &testHelper{
		ctx:      ctx,
		tprep:    tprep,
		dbCnfg:   dbCnfg,
		userIDs:  []uuid.UUID{userID1, userID2, userID3},
		eventIDs: []uuid.UUID{eventID1, eventID2, eventID3},
	}
}

func (th *testHelper) createTestTicketPurchase(num int, eventID uuid.UUID, userID uuid.UUID) *models.TicketPurchase {
	tp, err := models.NewTicketPurchase(
		uuid.New(),
		fmt.Sprintf("Customer %d", num),
		fmt.Sprintf("customer%d@example.com", num),
		time.Now().Add(time.Duration(num)*time.Hour),
		eventID,
		userID,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestTicketPurchase failed: %v", err))
	}
	return &tp
}

func (th *testHelper) createAndAddTicketPurchase(t *testing.T, num int, eventID uuid.UUID, userID uuid.UUID) *models.TicketPurchase {
	tp := th.createTestTicketPurchase(num, eventID, userID)

	err := th.tprep.Add(th.ctx, tp)
	require.NoError(t, err)

	return tp
}

func TestTicketPurchasesRep_GetTPurchasesOfUserID(t *testing.T) {
	th := setupTestHelper(t)

	userID := th.userIDs[0]
	eventID := th.eventIDs[0]
	tp1 := th.createAndAddTicketPurchase(t, 1, eventID, userID)
	tp2 := th.createAndAddTicketPurchase(t, 2, eventID, userID)
	otherUserID := th.userIDs[1]
	th.createAndAddTicketPurchase(t, 3, eventID, otherUserID)

	tests := []struct {
		name          string
		userID        uuid.UUID
		wantLen       int
		wantIDs       []uuid.UUID
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for user with no purchases",
			userID:        th.userIDs[2],
			wantLen:       0,
			expectedError: nil,
		},
		{
			name:    "Should return all ticket purchases for user",
			userID:  userID,
			wantLen: 2,
			wantIDs: []uuid.UUID{tp1.GetID(), tp2.GetID()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tps, err := th.tprep.GetTPurchasesOfUserID(th.ctx, tt.userID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Len(t, tps, tt.wantLen)

			if len(tt.wantIDs) > 0 {
				for i, expectedID := range tt.wantIDs {
					assert.Equal(t, expectedID, tps[i].GetID())
				}
			}
		})
	}
}

func TestTicketPurchasesRep_GetCntTPurchasesForEvent(t *testing.T) {
	th := setupTestHelper(t)

	eventID := th.eventIDs[0]
	userID := th.userIDs[0]

	// Add 3 tickets for the same event
	for i := 0; i < 3; i++ {
		tp := th.createTestTicketPurchase(i, eventID, userID)
		err := th.tprep.Add(th.ctx, tp)
		require.NoError(t, err)
	}

	// Add 1 ticket for another event
	otherEventID := th.eventIDs[1]
	tpOther := th.createTestTicketPurchase(4, otherEventID, userID)
	err := th.tprep.Add(th.ctx, tpOther)
	require.NoError(t, err)

	tests := []struct {
		name      string
		eventID   uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Should return correct count for event with purchases",
			eventID:   eventID,
			wantCount: 3,
		},
		{
			name:      "Should return 0 for event with no purchases",
			eventID:   th.eventIDs[2],
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := th.tprep.GetCntTPurchasesForEvent(th.ctx, tt.eventID)

			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, count)
		})
	}
}

func TestTicketPurchasesRep_Add(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should add new ticket purchase", func(t *testing.T) {
		tp := th.createTestTicketPurchase(1, th.eventIDs[0], th.userIDs[0])

		err := th.tprep.Add(th.ctx, tp)
		require.NoError(t, err)

		// Проверяем, что можно получить добавленную покупку
		tps, err := th.tprep.GetTPurchasesOfUserID(th.ctx, tp.GetUserID())
		require.NoError(t, err)
		assert.Len(t, tps, 1)
		assert.Equal(t, tp.GetID(), tps[0].GetID())
	})

	t.Run("Should return error for duplicate ticket purchase", func(t *testing.T) {
		tp := th.createAndAddTicketPurchase(t, 2, th.eventIDs[0], th.userIDs[0])

		// Пытаемся добавить тот же билет снова
		err := th.tprep.Add(th.ctx, tp)
		assert.Error(t, err)
	})

	t.Run("Should add ticket purchase without userID", func(t *testing.T) {
		tp := th.createTestTicketPurchase(3, th.eventIDs[0], uuid.Nil)

		err := th.tprep.Add(th.ctx, tp)
		require.NoError(t, err)

		// Проверяем, что билет добавлен, но не привязан к пользователю
		count, err := th.tprep.GetCntTPurchasesForEvent(th.ctx, tp.GetEventID())
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})
}

func TestTicketPurchasesRep_Ping(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should ping database successfully", func(t *testing.T) {
		err := th.tprep.Ping(th.ctx)
		require.NoError(t, err)
	})

	t.Run("Should return error after closing connection", func(t *testing.T) {
		th.tprep.Close()
		err := th.tprep.Ping(th.ctx)
		assert.Error(t, err)
	})
}
