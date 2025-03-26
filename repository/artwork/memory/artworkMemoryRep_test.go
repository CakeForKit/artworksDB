package artworkMemoryRep

import (
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/models"
	artworkRep "git.iu7.bmstu.ru/ped22u691/PPO.git/repository/artwork"
	"github.com/google/uuid"
)

func TestArtworkMemoryRep_GetByID(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	auth, err := models.NewAuthor("authorName", 1900, 2000)
	if err != nil {
		t.Fatal(err)
	}
	col, err := models.NewCollection("collectiontitle")
	if err != nil {
		t.Fatal(err)
	}
	aw, err := models.NewArtwork(
		"testTitle", 2000,
		&auth, &col,
		"20x20", "artwork material", "artwork technic",
	)
	if err != nil {
		t.Fatal(err)
	}

	var id uuid.UUID = aw.GetID()
	var rep ArtworkMemoryRep = ArtworkMemoryRep{
		artworks: map[uuid.UUID]models.Artwork{
			id: aw,
		},
	}

	var testCases []testCase = []testCase{
		{
			name:        "Artwork by ID",
			id:          id,
			expectedErr: nil,
		}, {
			name:        "No artwork by ID",
			id:          uuid.New(),
			expectedErr: artworkRep.ErrArtworkNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rep.Get(tc.id)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
