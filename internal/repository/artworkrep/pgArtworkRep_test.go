package artworkrep_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	th     *testHelper
	pgOnce sync.Once
)

// testHelper содержит общие методы для тестов
type testHelper struct {
	ctx          context.Context
	arep         *artworkrep.PgArtworkRep
	dbCnfg       *cnfg.DatebaseConfig
	pgTestConfig *cnfg.PostgresTestConfig
	pgCreds      *cnfg.PostgresCredentials
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {

		dbCnfg := cnfg.GetTestDatebaseConfig()
		pgTestConfig := cnfg.GetPgTestConfig()

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		arep, err := artworkrep.NewPgArtworkRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:          ctx,
			arep:         arep,
			dbCnfg:       dbCnfg,
			pgTestConfig: pgTestConfig,
			pgCreds:      &pgCreds,
		}
	})
	err := pgtest.MigrateUp(ctx, th.pgTestConfig, th.pgCreds)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := pgtest.MigrateDown(ctx, th.pgTestConfig, th.pgCreds)
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
	err := th.arep.Add(context.Background(), artwork)
	require.NoError(t, err)

	return artwork, author, collection
}

func TestArtworkRep_GetAll(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() []*models.Artwork
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "Should return empty list for empty DB",
			setup:         func() []*models.Artwork { return nil },
			wantLen:       0,
			expectedError: artworkrep.ErrArtworkNotFound,
		},
		{
			name: "Should return all artworks",
			setup: func() []*models.Artwork {
				artworks := make([]*models.Artwork, 3)
				for i := range artworks {
					artwork, _, _ := th.createAndAddArtwork(t, i)
					artworks[i] = artwork
				}
				return artworks
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedArtworks := tt.setup()

			artworks, err := th.arep.GetAll(th.ctx)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Len(t, artworks, tt.wantLen)

			if len(expectedArtworks) > 0 {
				for i, expected := range expectedArtworks {
					assert.Equal(t, expected.GetID(), artworks[i].GetID())
					assert.Equal(t, expected.GetTitle(), artworks[i].GetTitle())
				}
			}
		})
	}
}

func TestArtworkRep_GetByID(t *testing.T) {
	th := setupTestHelper(t)

	artwork, _, _ := th.createAndAddArtwork(t, 1)
	nonExistentID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *models.Artwork
		wantErr error
	}{
		{
			name: "Should return artwork by ID",
			id:   artwork.GetID(),
			want: artwork,
		},
		{
			name:    "Should return error for non-existent ID",
			id:      nonExistentID,
			wantErr: artworkrep.ErrArtworkNotFound,
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
			assert.Equal(t, tt.want.GetTitle(), got.GetTitle())
		})
	}
}

func TestArtworkRep_GetByTitle(t *testing.T) {
	th := setupTestHelper(t)

	artwork, _, _ := th.createAndAddArtwork(t, 1)

	t.Run("Should return artwork by title", func(t *testing.T) {
		artworks, err := th.arep.GetByTitle(th.ctx, artwork.GetTitle())
		require.NoError(t, err)
		assert.Len(t, artworks, 1)
		assert.Equal(t, artwork.GetID(), artworks[0].GetID())
	})

	t.Run("Should return error for non-existent title", func(t *testing.T) {
		res, err := th.arep.GetByTitle(th.ctx, "non_existent_title")
		assert.ErrorIs(t, err, artworkrep.ErrArtworkNotFound)
		assert.Nil(t, res)
	})
}

func TestArtworkRep_GetByAuthor(t *testing.T) {
	th := setupTestHelper(t)

	artwork, author, _ := th.createAndAddArtwork(t, 1)

	t.Run("Should return artworks by author", func(t *testing.T) {
		artworks, err := th.arep.GetByAuthor(th.ctx, author)
		require.NoError(t, err)
		assert.Len(t, artworks, 1)
		assert.Equal(t, artwork.GetID(), artworks[0].GetID())
	})

	t.Run("Should return error for non-existent author", func(t *testing.T) {
		nonExistentAuthor := &models.Author{}
		res, err := th.arep.GetByAuthor(th.ctx, nonExistentAuthor)
		assert.ErrorIs(t, err, artworkrep.ErrArtworkNotFound)
		assert.Nil(t, res)
	})
}

func TestArtworkRep_GetByCreationTime(t *testing.T) {
	th := setupTestHelper(t)

	artwork, _, _ := th.createAndAddArtwork(t, 1)
	year := artwork.GetCreationYear()

	tests := []struct {
		name     string
		yearBeg  int
		yearEnd  int
		wantLen  int
		wantErr  bool
		wantYear int
	}{
		{
			name:     "Should return artworks in time range",
			yearBeg:  year - 1,
			yearEnd:  year + 1,
			wantLen:  1,
			wantYear: year,
		},
		{
			name:    "Should return empty for out of range",
			yearBeg: year + 10,
			yearEnd: year + 20,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artworks, err := th.arep.GetByCreationTime(th.ctx, tt.yearBeg, tt.yearEnd)

			if tt.wantLen == 0 {
				assert.ErrorIs(t, err, artworkrep.ErrArtworkNotFound)
				return
			}

			require.NoError(t, err)
			assert.Len(t, artworks, tt.wantLen)
			if tt.wantLen > 0 {
				assert.Equal(t, tt.wantYear, artworks[0].GetCreationYear())
			}
		})
	}
}

func TestArtworkRep_Add(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should add new artwork", func(t *testing.T) {
		author := th.createTestAuthor(1)
		collection := th.createTestCollection(1)
		artwork := th.createTestArtwork(1, author, collection)

		err := th.arep.Add(th.ctx, artwork)
		require.NoError(t, err)

		// Проверяем, что artwork действительно добавлен
		got, err := th.arep.GetByID(th.ctx, artwork.GetID())
		require.NoError(t, err)
		assert.Equal(t, artwork.GetID(), got.GetID())
	})

	t.Run("Should return error for duplicate artwork", func(t *testing.T) {
		artwork, _, _ := th.createAndAddArtwork(t, 2)

		// Пытаемся добавить тот же artwork снова
		err := th.arep.Add(th.ctx, artwork)
		assert.Error(t, err)
	})
}

func TestArtworkRep_Delete(t *testing.T) {
	th := setupTestHelper(t)

	t.Run("Should delete existing artwork", func(t *testing.T) {
		artwork, _, _ := th.createAndAddArtwork(t, 1)

		err := th.arep.Delete(th.ctx, artwork.GetID())
		require.NoError(t, err)

		// Проверяем, что artwork удален
		_, err = th.arep.GetByID(th.ctx, artwork.GetID())
		assert.ErrorIs(t, err, artworkrep.ErrArtworkNotFound)
	})

	t.Run("Should return error for non-existent artwork", func(t *testing.T) {
		err := th.arep.Delete(th.ctx, uuid.New())
		assert.ErrorIs(t, err, artworkrep.ErrRowsAffected)
	})
}

func TestArtworkRep_Update(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name       string
		updateFunc func(*models.Artwork) (*models.Artwork, error)
		wantCheck  func(*testing.T, *models.Artwork)
		wantErr    bool
	}{
		{
			name: "Should update artwork",
			updateFunc: func(a *models.Artwork) (*models.Artwork, error) {
				newArtwork, err := models.NewArtwork(
					a.GetID(),
					"Updated Title",
					"Updated Technic",
					"Updated Material",
					"Updated Size",
					a.GetCreationYear()+1,
					a.GetAuthor(),
					a.GetCollection(),
				)
				return &newArtwork, err
			},
			wantCheck: func(t *testing.T, a *models.Artwork) {
				assert.Equal(t, "Updated Title", a.GetTitle())
				assert.Equal(t, "Updated Technic", a.GetTechnic())
			},
		},
		{
			name: "Should return error from updateFunc",
			updateFunc: func(a *models.Artwork) (*models.Artwork, error) {
				return nil, errors.New("update error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artwork, _, _ := th.createAndAddArtwork(t, 1)

			updated, err := th.arep.Update(th.ctx, artwork.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.wantCheck(t, updated)

			// Проверяем, что изменения сохранились в БД
			dbArtwork, err := th.arep.GetByID(th.ctx, artwork.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbArtwork)
		})
	}
}
