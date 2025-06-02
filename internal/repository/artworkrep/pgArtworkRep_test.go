package artworkrep_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
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
	ctx       context.Context
	arep      *artworkrep.PgArtworkRep
	authorRep *authorrep.PgAuthorRep
	colRep    *collectionrep.PgCollectionRep
	dbCnfg    *cnfg.DatebaseConfig
	pgCreds   *cnfg.DatebaseCredentials
}

func setupTestHelper(t *testing.T) *testHelper {
	ctx := context.Background()
	pgOnce.Do(func() {
		dbCnfg := cnfg.GetTestDatebaseConfig()

		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		require.NoError(t, err)

		arep, err := artworkrep.NewPgArtworkRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)
		collectionRep, err := collectionrep.NewPgCollectionRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)
		authorRep, err := authorrep.NewPgAuthorRep(ctx, &pgCreds, dbCnfg)
		require.NoError(t, err)

		th = &testHelper{
			ctx:       ctx,
			arep:      arep,
			dbCnfg:    dbCnfg,
			pgCreds:   &pgCreds,
			colRep:    collectionRep,
			authorRep: authorRep,
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

func (th *testHelper) deleteArtwork(t *testing.T, artworks []*models.Artwork) {
	ctx := context.Background()
	for _, artwork := range artworks {
		err := th.arep.Delete(ctx, artwork.GetID())
		require.NoError(t, err)
		err = th.authorRep.Delete(ctx, artwork.GetAuthor().GetID())
		require.NoError(t, err)
		err = th.colRep.DeleteCollection(ctx, artwork.GetCollection().GetID())
		require.NoError(t, err)
	}
}

func TestPgArtworkRep_GetAllArtworks(t *testing.T) {
	th := setupTestHelper(t)

	tests := []struct {
		name          string
		setup         func() ([]*models.Artwork, []*models.Artwork)
		down          func([]*models.Artwork)
		filter        *jsonreqresp.ArtworkFilter
		sort          *jsonreqresp.ArtworkSortOps
		wantLen       int
		wantErr       bool
		expectedError error
	}{
		{
			name: "Should return empty list for empty DB",
			setup: func() ([]*models.Artwork, []*models.Artwork) {
				return nil, nil
			},
			down:          func(a []*models.Artwork) {},
			filter:        &jsonreqresp.ArtworkFilter{},
			sort:          &jsonreqresp.ArtworkSortOps{},
			wantLen:       0,
			expectedError: nil,
		},
		{
			name: "Should return all artworks",
			setup: func() ([]*models.Artwork, []*models.Artwork) {
				artworks := make([]*models.Artwork, 3)
				for i := range artworks {
					artwork, _, _ := th.createAndAddArtwork(t, i)
					artworks[i] = artwork
				}
				return artworks, artworks
			},
			down:    func(a []*models.Artwork) { th.deleteArtwork(t, a) },
			filter:  &jsonreqresp.ArtworkFilter{},
			sort:    &jsonreqresp.ArtworkSortOps{},
			wantLen: 3,
		},
		{
			name: "Should filter by title",
			setup: func() ([]*models.Artwork, []*models.Artwork) {
				artworks := make([]*models.Artwork, 3)
				for i := range artworks {
					artwork, _, _ := th.createAndAddArtwork(t, i)
					artworks[i] = artwork
				}
				return artworks[:1], artworks // Expect only first artwork
			},
			down: func(a []*models.Artwork) { th.deleteArtwork(t, a) },
			filter: &jsonreqresp.ArtworkFilter{
				Title: "Artwork 0",
			},
			sort:    &jsonreqresp.ArtworkSortOps{},
			wantLen: 1,
		},
		{
			name: "Should sort by title",
			setup: func() ([]*models.Artwork, []*models.Artwork) {
				artworks := make([]*models.Artwork, 3)
				for i := range artworks {
					artwork, _, _ := th.createAndAddArtwork(t, i)
					artworks[i] = artwork
				}
				return artworks, artworks
			},
			down:   func(a []*models.Artwork) { th.deleteArtwork(t, a) },
			filter: &jsonreqresp.ArtworkFilter{},
			sort: &jsonreqresp.ArtworkSortOps{
				Field:     jsonreqresp.TitleSortFieldArtwork,
				Direction: "ASC",
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedArtworks, AllArtworks := tt.setup()

			artworks, err := th.arep.GetAllArtworks(th.ctx, tt.filter, tt.sort)

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
			tt.down(AllArtworks)
		})
	}
}

func TestPgArtworkRep_GetByID(t *testing.T) {
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

func TestPgArtworkRep_Update(t *testing.T) {
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

			err := th.arep.Update(th.ctx, artwork.GetID(), tt.updateFunc)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			dbArtwork, err := th.arep.GetByID(th.ctx, artwork.GetID())
			require.NoError(t, err)
			tt.wantCheck(t, dbArtwork)
		})
	}
}
