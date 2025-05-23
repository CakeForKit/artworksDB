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
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
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
	ctx        context.Context
	erep       *eventrep.PgEventRep
	arep       *artworkrep.PgArtworkRep
	authorRep  *authorrep.PgAuthorRep
	colRep     *collectionrep.PgCollectionRep
	dbCnfg     *cnfg.DatebaseConfig
	pgCreds    *cnfg.PostgresCredentials
	employeeID uuid.UUID
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

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		erep, err := eventrep.NewPgEventRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)
		arep, err := artworkrep.NewPgArtworkRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)
		collectionRep, err := collectionrep.NewPgCollectionRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)
		authorRep, err := authorrep.NewPgAuthorRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:        ctx,
			erep:       erep,
			arep:       arep,
			colRep:     collectionRep,
			authorRep:  authorRep,
			employeeID: uuid.New(),
			dbCnfg:     dbCnfg,
			pgCreds:    &pgCreds,
		}
	})
	pgTestConfig := cnfg.GetPgTestConfig()
	err := pgtest.MigrateUp(ctx, pgTestConfig.MigrationDir, th.pgCreds)
	require.NoError(t, err)
	addEmployee(t, ctx, th.employeeID, th.pgCreds, th.dbCnfg)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx, pgTestConfig.MigrationDir, th.pgCreds)
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
		true,
		nil,
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

func (th *testHelper) deleteEvent(t *testing.T, events []*models.Event) {
	ctx := context.Background()
	for _, event := range events {
		err := th.erep.RealDelete(ctx, event.GetID())
		require.NoError(t, err)
	}
}

func (th *testHelper) createTestAuthor(num int) *models.Author {
	author, err := models.NewAuthor(
		uuid.New(),
		fmt.Sprintf("Author %d", num),
		1900+num,
		1950+num,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestAuthor failed: %v", err))
	}
	return &author
}

func (th *testHelper) createTestCollection(num int) *models.Collection {
	collection, err := models.NewCollection(
		uuid.New(),
		fmt.Sprintf("Collection %d", num),
	)
	if err != nil {
		panic(fmt.Sprintf("createTestCollection failed: %v", err))
	}
	return &collection
}

func (th *testHelper) createTestArtwork(num int, author *models.Author, collection *models.Collection) *models.Artwork {
	artwork, err := models.NewArtwork(
		uuid.New(),
		fmt.Sprintf("Artwork %d", num),
		"Oil on canvas",
		"Canvas",
		"100x100 cm",
		1920+num,
		author,
		collection,
	)
	if err != nil {
		panic(fmt.Sprintf("createTestArtwork failed: %v", err))
	}
	return &artwork
}

func (th *testHelper) createAndAddArtwork(t *testing.T, num int) (*models.Artwork, *models.Author, *models.Collection) {
	author := th.createTestAuthor(num)
	collection := th.createTestCollection(num)
	artwork := th.createTestArtwork(num, author, collection)

	// Добавляем автора и коллекцию сначала
	ctx := context.Background()
	err := th.authorRep.Add(ctx, author)
	require.NoError(t, err)
	err = th.colRep.AddCollection(ctx, collection)
	require.NoError(t, err)
	err = th.arep.Add(ctx, artwork)
	require.NoError(t, err)

	return artwork, author, collection
}

func TestEventRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		filter        *jsonreqresp.EventFilter
		setup         func() ([]*models.Event, []*models.Event)
		down          func([]*models.Event)
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			filter:        &jsonreqresp.EventFilter{},
			setup:         func() ([]*models.Event, []*models.Event) { return nil, nil },
			down:          func(a []*models.Event) {},
			wantLen:       0,
			expectedError: nil,
		},
		{
			name:   "Should return all events",
			filter: &jsonreqresp.EventFilter{},
			setup: func() ([]*models.Event, []*models.Event) {
				events := make([]*models.Event, 3)
				for i := range events {
					events[i] = th.createAndAddEvent(t, i)
				}
				return events, events
			},
			down:    func(a []*models.Event) { th.deleteEvent(t, a) },
			wantLen: 3,
		},
		{
			name: "Should filter by title",
			filter: &jsonreqresp.EventFilter{
				Title: "Event 1",
			},
			setup: func() ([]*models.Event, []*models.Event) {
				events := make([]*models.Event, 3)
				for i := range events {
					events[i] = th.createAndAddEvent(t, i)
				}
				return []*models.Event{events[1]}, events
			},
			down:    func(a []*models.Event) { th.deleteEvent(t, a) },
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedEvents, AllEvents := tt.setup()
			events, err := th.erep.GetAll(th.ctx, tt.filter)

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
			tt.down(AllEvents)
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
			assert.Equal(t, tt.want.GetAddress(), got.GetAddress())
			assert.Equal(t, tt.want.GetEmployeeID(), got.GetEmployeeID())
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
			filter := &jsonreqresp.EventFilter{
				DateBegin: tt.dateBeg,
				DateEnd:   tt.dateEnd,
			}
			events, err := th.erep.GetAll(th.ctx, filter)

			if tt.wantLen == 0 {
				assert.ErrorIs(t, err, nil)
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
					true,
					nil,
				)
				return &newEvent, err
			},
			wantCheck: func(t *testing.T, e *models.Event) {
				assert.Equal(t, "Updated Title", e.GetTitle())
				assert.Equal(t, "Updated Address", e.GetAddress())
				assert.Equal(t, 200, e.GetTicketCount())
				assert.False(t, e.GetAccess())
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

			err := th.erep.Update(th.ctx, event.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			updated, err := th.erep.GetByID(th.ctx, event.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, updated)
		})
	}
}

func TestEventRep_ArtworkOperations(t *testing.T) {
	th := setupTestHelper(t)
	event := th.createAndAddEvent(t, 1)
	art, _, _ := th.createAndAddArtwork(t, 1)
	artworkID := art.GetID()

	t.Run("Add artwork to event", func(t *testing.T) {
		err := th.erep.AddArtworksToEvent(th.ctx, event.GetID(), uuid.UUIDs{artworkID})
		require.NoError(t, err)

		got, err := th.erep.GetByID(th.ctx, event.GetID())
		require.NoError(t, err)
		assert.Contains(t, got.GetArtworkIDs(), artworkID)
	})

	t.Run("Delete artwork from event", func(t *testing.T) {
		err := th.erep.DeleteArtworkFromEvent(th.ctx, event.GetID(), artworkID)
		require.NoError(t, err)

		got, err := th.erep.GetByID(th.ctx, event.GetID())
		require.NoError(t, err)
		assert.NotContains(t, got.GetArtworkIDs(), artworkID)
	})
}
