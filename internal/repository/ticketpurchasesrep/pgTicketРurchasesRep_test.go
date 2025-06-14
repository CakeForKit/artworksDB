package ticketpurchasesrep_test

import (
	"context"
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
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	th     *testHelper
	pgOnce sync.Once
)

type testHelper struct {
	ctx        context.Context
	tprep      *ticketpurchasesrep.PgTicketPurchasesRep
	dbCnfg     *cnfg.DatebaseConfig
	pgCreds    *cnfg.DatebaseCredentials
	userIDs    []uuid.UUID
	eventIDs   []uuid.UUID
	employeeID uuid.UUID
}

func addUser(t *testing.T, ctx context.Context, userID uuid.UUID, pgCreds *cnfg.DatebaseCredentials, dbCnfg *cnfg.DatebaseConfig, num int) {
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
}

func addEmployee(t *testing.T, ctx context.Context, employeeID uuid.UUID, pgCreds *cnfg.DatebaseCredentials, dbCnfg *cnfg.DatebaseConfig) {
	admin, err := models.NewAdmin(
		uuid.New(),
		"admin",
		"adminL",
		"hashPadmin",
		time.Now().UTC().Truncate(time.Microsecond),
		true,
	)
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

func addEvent(t *testing.T, ctx context.Context, eventID uuid.UUID, pgCreds *cnfg.DatebaseCredentials, dbCnfg *cnfg.DatebaseConfig, employeeID uuid.UUID, num int) {
	event, err := models.NewEvent(
		eventID,
		fmt.Sprintf("Event %d", num),
		time.Now().Add(time.Hour),
		time.Now().Add(2*time.Hour),
		fmt.Sprintf("Address %d", num),
		true,
		employeeID,
		100+num,
		true,
		nil,
	)
	require.NoError(t, err)
	erep, err := eventrep.NewPgEventRep(ctx, pgCreds, dbCnfg)
	require.NoError(t, err)
	err = erep.Add(ctx, &event)
	require.NoError(t, err)
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {
		dbCnfg := cnfg.GetTestDatebaseConfig()

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		tprep, err := ticketpurchasesrep.NewPgTicketPurchasesRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:        ctx,
			tprep:      tprep,
			dbCnfg:     dbCnfg,
			pgCreds:    &pgCreds,
			userIDs:    []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			eventIDs:   []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			employeeID: uuid.New(),
		}
	})
	pgTestConfig := cnfg.GetPgTestConfig()
	err := pgtest.MigrateUp(ctx, pgTestConfig.MigrationDir, th.pgCreds)
	require.NoError(t, err)

	addEmployee(t, ctx, th.employeeID, th.pgCreds, th.dbCnfg)
	addUser(t, ctx, th.userIDs[0], th.pgCreds, th.dbCnfg, 1)
	addUser(t, ctx, th.userIDs[1], th.pgCreds, th.dbCnfg, 2)
	addUser(t, ctx, th.userIDs[2], th.pgCreds, th.dbCnfg, 3)
	addEvent(t, ctx, th.eventIDs[0], th.pgCreds, th.dbCnfg, th.employeeID, 1)
	addEvent(t, ctx, th.eventIDs[1], th.pgCreds, th.dbCnfg, th.employeeID, 2)
	addEvent(t, ctx, th.eventIDs[2], th.pgCreds, th.dbCnfg, th.employeeID, 3)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx, pgTestConfig.MigrationDir, th.pgCreds)
		require.NoError(t, err)
	})

	return th
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
					assert.Equal(t, fmt.Sprintf("Customer %d", i+1), tps[i].GetCustomerName())
					assert.Equal(t, fmt.Sprintf("customer%d@example.com", i+1), tps[i].GetCustomerEmail())
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
