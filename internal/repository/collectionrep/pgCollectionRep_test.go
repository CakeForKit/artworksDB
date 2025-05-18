package collectionrep_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	th     *testHelper
	pgOnce sync.Once
)

type testHelper struct {
	ctx     context.Context
	crep    *collectionrep.PgCollectionRep
	dbCnfg  *cnfg.DatebaseConfig
	pgCreds *cnfg.PostgresCredentials
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {
		dbCnfg := cnfg.GetTestDatebaseConfig()

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		crep, err := collectionrep.NewPgCollectionRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:     ctx,
			crep:    crep,
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

func (th *testHelper) createAndAddCollection(t *testing.T, num int) *models.Collection {
	collection := th.createTestCollection(num)
	err := th.crep.AddCollection(th.ctx, collection)
	require.NoError(t, err)
	return collection
}

func TestPgCollectionRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() []*models.Collection
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			setup:         func() []*models.Collection { return nil },
			wantLen:       0,
			expectedError: nil,
		},
		{
			name: "Should return all collections",
			setup: func() []*models.Collection {
				collections := make([]*models.Collection, 3)
				for i := range collections {
					collections[i] = th.createAndAddCollection(t, i)
				}
				return collections
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedCollections := tt.setup()

			collections, err := th.crep.GetAllCollections(th.ctx)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Len(t, collections, tt.wantLen)

			if len(expectedCollections) > 0 {
				for i, expected := range expectedCollections {
					assert.Equal(t, expected.GetID(), collections[i].GetID())
					assert.Equal(t, expected.GetTitle(), collections[i].GetTitle())
				}
			}
		})
	}
}

func TestPgCollectionRep_GetByID(t *testing.T) {
	th := setupTestHelper(t)

	collection := th.createAndAddCollection(t, 1)
	nonExistentID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *models.Collection
		wantErr error
	}{
		{
			name: "Should return collection by ID",
			id:   collection.GetID(),
			want: collection,
		},
		{
			name:    "Should return error for non-existent ID",
			id:      nonExistentID,
			wantErr: collectionrep.ErrCollectionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := th.crep.GetCollectionByID(th.ctx, tt.id)

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

func TestPgCollectionRep_Update(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		updateFunc func(*models.Collection) (*models.Collection, error)
		wantCheck  func(*testing.T, *models.Collection)
		wantErr    bool
	}{
		{
			name: "Should update collection",
			updateFunc: func(c *models.Collection) (*models.Collection, error) {
				newCollection, err := models.NewCollection(
					c.GetID(),
					"Updated Title",
				)
				return &newCollection, err
			},
			wantCheck: func(t *testing.T, c *models.Collection) {
				assert.Equal(t, "Updated Title", c.GetTitle())
			},
		},
		{
			name: "Should return error from updateFunc",
			updateFunc: func(c *models.Collection) (*models.Collection, error) {
				return nil, errors.New("update error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection := th.createAndAddCollection(t, 1)

			err := th.crep.UpdateCollection(th.ctx, collection.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify changes were persisted
			dbCollection, err := th.crep.GetCollectionByID(th.ctx, collection.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbCollection)
		})
	}
}
