package authorrep_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
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
	arep    *authorrep.PgAuthorRep
	dbCnfg  *cnfg.DatebaseConfig
	pgCreds *cnfg.PostgresCredentials
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {
		dbCnfg := cnfg.GetTestDatebaseConfig()

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		arep, err := authorrep.NewPgAuthorRep(ctx, &pgCreds, dbCnfg)
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

func (th *testHelper) createAndAddAuthor(t *testing.T, num int) *models.Author {
	author := th.createTestAuthor(num)
	err := th.arep.Add(th.ctx, author)
	require.NoError(t, err)
	return author
}

func TestPgAuthorRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() []*models.Author
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			setup:         func() []*models.Author { return nil },
			wantLen:       0,
			expectedError: nil,
		},
		{
			name: "Should return all authors",
			setup: func() []*models.Author {
				authors := make([]*models.Author, 3)
				for i := range authors {
					authors[i] = th.createAndAddAuthor(t, i)
				}
				return authors
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedAuthors := tt.setup()

			authors, err := th.arep.GetAll(th.ctx)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Len(t, authors, tt.wantLen)

			if len(expectedAuthors) > 0 {
				for i, expected := range expectedAuthors {
					assert.Equal(t, expected.GetID(), authors[i].GetID())
					assert.Equal(t, expected.GetName(), authors[i].GetName())
				}
			}
		})
	}
}

func TestPgAuthorRep_GetByID(t *testing.T) {
	th := setupTestHelper(t)

	author := th.createAndAddAuthor(t, 1)
	nonExistentID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *models.Author
		wantErr error
	}{
		{
			name: "Should return author by ID",
			id:   author.GetID(),
			want: author,
		},
		{
			name:    "Should return error for non-existent ID",
			id:      nonExistentID,
			wantErr: authorrep.ErrAuthorNotFound,
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
			assert.Equal(t, tt.want.GetName(), got.GetName())
		})
	}
}

func TestPgAuthorRep_Update(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		updateFunc func(*models.Author) (*models.Author, error)
		wantCheck  func(*testing.T, *models.Author)
		wantErr    bool
	}{
		{
			name: "Should update author",
			updateFunc: func(a *models.Author) (*models.Author, error) {
				newAuthor, err := models.NewAuthor(
					a.GetID(),
					"Updated Name",
					a.GetBirthYear()+1,
					a.GetDeathYear()+10,
				)
				return &newAuthor, err
			},
			wantCheck: func(t *testing.T, a *models.Author) {
				assert.Equal(t, "Updated Name", a.GetName())
				assert.Equal(t, 1902, a.GetBirthYear())
			},
		},
		{
			name: "Should return error from updateFunc",
			updateFunc: func(a *models.Author) (*models.Author, error) {
				return nil, errors.New("update error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			author := th.createAndAddAuthor(t, 1)

			err := th.arep.Update(th.ctx, author.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify changes were persisted
			dbAuthor, err := th.arep.GetByID(th.ctx, author.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbAuthor)
		})
	}
}

func TestPgAuthorRep_HasArtworks(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should return false for author without artworks", func(t *testing.T) {
		author := th.createAndAddAuthor(t, 1)

		hasArtworks, err := th.arep.HasArtworks(th.ctx, author.GetID())
		require.NoError(t, err)
		assert.False(t, hasArtworks)
	})

	// Note: This test would require setting up artwork relationships
	// which would need artwork repository and proper database setup
	t.Run("Should return true for author with artworks", func(t *testing.T) {
		t.Skip("Requires artwork relationship setup")
	})
}
