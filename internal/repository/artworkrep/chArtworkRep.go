package artworkrep

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
)

type CHArtworkRep struct {
	db *sql.DB
}

var (
	chInstance *CHArtworkRep
	chOnce     sync.Once
)

func NewCHArtworkRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHArtworkRep, error) {
	var resErr error
	chOnce.Do(func() {
		conn := clickhouse.OpenDB(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", chCreds.Host, chCreds.Port)},
			Auth: clickhouse.Auth{
				Database: chCreds.DbName,
				Username: chCreds.Username,
				Password: chCreds.Password,
			},
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			},
		})

		if err := conn.PingContext(ctx); err != nil {
			resErr = fmt.Errorf("NewCHArtworkRep: %w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHArtworkRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHArtworkRep) parseArtworksRows(rows *sql.Rows) ([]*models.Artwork, error) {
	var resArtworks []*models.Artwork
	for rows.Next() {
		var id, authorID, collectionID uuid.UUID
		var title, authorName, collectionTitle, size, material, technic string
		var creationYear, authorBirthYear int32
		var authorDeathYear sql.NullInt32

		if err := rows.Scan(&id, &title, &technic, &material, &size, &creationYear,
			&authorID, &authorName, &authorBirthYear, &authorDeathYear,
			&collectionID, &collectionTitle); err != nil {
			return nil, fmt.Errorf("parseArtworksRows: scan error: %v", err)
		}

		deathYear := 0
		if authorDeathYear.Valid {
			deathYear = int(authorDeathYear.Int32)
		}

		author, err := models.NewAuthor(authorID, authorName, int(authorBirthYear), deathYear)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %v", err)
		}

		collection, err := models.NewCollection(collectionID, collectionTitle)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %v", err)
		}

		artwork, err := models.NewArtwork(id, title, technic, material, size, int(creationYear), &author, &collection)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %w: %v", models.ErrValidateArtwork, err)
		}
		resArtworks = append(resArtworks, &artwork)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("parseArtworksRows: rows iteration error: %v", err)
	}
	return resArtworks, nil
}

func (ch *CHArtworkRep) buildFilterConditions(filterOps *jsonreqresp.ArtworkFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filterOps.Title != "" {
		conditions = append(conditions, "Artworks.title LIKE ?")
		args = append(args, "%"+filterOps.Title+"%")
	}
	if filterOps.AuthorName != "" {
		conditions = append(conditions, "Author.name LIKE ?")
		args = append(args, "%"+filterOps.AuthorName+"%")
	}
	if filterOps.Collection != "" {
		conditions = append(conditions, "Collection.title LIKE ?")
		args = append(args, "%"+filterOps.Collection+"%")
	}
	if filterOps.EventID != uuid.Nil {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM Artwork_event ae WHERE Artworks.id = ae.artworkID AND ae.eventID = ?)")
		args = append(args, filterOps.EventID)
	}

	if len(conditions) == 0 {
		return "", nil
	}

	return "WHERE " + joinConditions(conditions, " AND "), args
}

func joinConditions(conditions []string, sep string) string {
	if len(conditions) == 0 {
		return ""
	}
	if len(conditions) == 1 {
		return conditions[0]
	}

	result := conditions[0]
	for _, cond := range conditions[1:] {
		result += sep + cond
	}
	return result
}

func (ch *CHArtworkRep) buildSortClause(sortOps *jsonreqresp.ArtworkSortOps) string {
	if sortOps == nil {
		return ""
	}

	switch sortOps.Field {
	case jsonreqresp.TitleSortFieldArtwork:
		return "ORDER BY Artworks.title " + sortOps.Direction
	case jsonreqresp.AuthorNameSortFieldArtwork:
		return "ORDER BY Author.name " + sortOps.Direction
	case jsonreqresp.CreationYearSortFieldArtwork:
		return "ORDER BY Artworks.creationYear " + sortOps.Direction
	default:
		return ""
	}
}

func (ch *CHArtworkRep) GetAllArtworks(ctx context.Context, filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error) {
	baseQuery := `
		SELECT 
			Artworks.id, Artworks.title, Artworks.technic, Artworks.material,
			Artworks.size, Artworks.creationYear,
			Author.id, Author.name, Author.birthYear, Author.deathYear,
			Collection.id, Collection.title
		FROM Artworks
		JOIN Author ON Artworks.authorID = Author.id
		JOIN Collection ON Artworks.collectionID = Collection.id`

	filterClause, filterArgs := ch.buildFilterConditions(filterOps)
	sortClause := ch.buildSortClause(sortOps)

	query := baseQuery
	if filterClause != "" {
		query += " " + filterClause
	}
	if sortClause != "" {
		query += " " + sortClause
	}

	rows, err := ch.db.QueryContext(ctx, query, filterArgs...)
	if err != nil {
		return nil, fmt.Errorf("CHArtworkRep.GetAllArtworks: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := ch.parseArtworksRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHArtworkRep.GetAllArtworks: %w", err)
	}
	return arts, nil
}

func (ch *CHArtworkRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Artwork, error) {
	query := `
		SELECT 
			Artworks.id, Artworks.title, Artworks.technic, Artworks.material,
			Artworks.size, Artworks.creationYear,
			Author.id, Author.name, Author.birthYear, Author.deathYear,
			Collection.id, Collection.title
		FROM Artworks
		JOIN Author ON Artworks.authorID = Author.id
		JOIN Collection ON Artworks.collectionID = Collection.id
		WHERE Artworks.id = ?`

	rows, err := ch.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("CHArtworkRep.GetByID: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := ch.parseArtworksRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHArtworkRep.GetByID: %w", err)
	}

	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("CHArtworkRep.GetByID: %w", ErrExpectedOneArtwork)
	}
	return arts[0], nil
}

func (ch *CHArtworkRep) execChangeQuery(ctx context.Context, query string, args ...interface{}) error {
	result, err := ch.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHArtworkRep) Add(ctx context.Context, a *models.Artwork) error {
	query := `
		INSERT INTO Artworks 
		(id, title, technic, material, size, creationYear, authorID, collectionID) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	err := ch.execChangeQuery(ctx, query,
		a.GetID(),
		a.GetTitle(),
		a.GetTechnic(),
		a.GetMaterial(),
		a.GetSize(),
		a.GetCreationYear(),
		a.GetAuthor().GetID(),
		a.GetCollection().GetID(),
	)
	if err != nil {
		return fmt.Errorf("CHArtworkRep.Add: %w", err)
	}
	return nil
}

func (ch *CHArtworkRep) Delete(ctx context.Context, idArt uuid.UUID) error {
	query := "ALTER TABLE Artworks DELETE WHERE id = ?"
	err := ch.execChangeQuery(ctx, query, idArt)
	if err != nil {
		return fmt.Errorf("CHArtworkRep.Delete: %w", err)
	}
	return nil
}

func (ch *CHArtworkRep) Update(ctx context.Context,
	idArt uuid.UUID,
	funcUpdate func(*models.Artwork) (*models.Artwork, error),
) error {
	art, err := ch.GetByID(ctx, idArt)
	if err != nil {
		return fmt.Errorf("CHArtworkRep.Update: %w", err)
	}

	updatedArtwork, err := funcUpdate(art)
	if err != nil {
		return fmt.Errorf("CHArtworkRep.Update: %w", ErrUpdateArtwork)
	}

	query := `
		ALTER TABLE Artworks UPDATE 
		title = ?, 
		material = ?, 
		technic = ?, 
		size = ?, 
		creationYear = ?, 
		authorID = ?, 
		collectionID = ? 
		WHERE id = ?`

	err = ch.execChangeQuery(ctx, query,
		updatedArtwork.GetTitle(),
		updatedArtwork.GetMaterial(),
		updatedArtwork.GetTechnic(),
		updatedArtwork.GetSize(),
		updatedArtwork.GetCreationYear(),
		updatedArtwork.GetAuthor().GetID(),
		updatedArtwork.GetCollection().GetID(),
		idArt,
	)
	if err != nil {
		return fmt.Errorf("CHArtworkRep.Update: %w", err)
	}
	return nil
}

func (ch *CHArtworkRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHArtworkRep) Close() {
	ch.db.Close()
}
