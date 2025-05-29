package artworkrep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgArtworkRep struct {
	db *sql.DB
}

var (
	pgInstance *PgArtworkRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect        = errors.New("open connect failed")
	ErrPing               = errors.New("ping failed")
	ErrQueryBuilds        = errors.New("query build failed")
	ErrQueryExec          = errors.New("query execution failed")
	ErrExpectedOneArtwork = errors.New("expected one artwork")
	ErrRowsAffected       = errors.New("no rows affected")
)

func NewPgArtworkRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (*PgArtworkRep, error) {
	var resErr error
	pgOnce.Do(func() {
		// connStr := "postgres://puser:ppassword@postgres_artworks:5432/artworks"
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			resErr = fmt.Errorf("NewPgArtworkRep: %w: %v", ErrOpenConnect, err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			resErr = fmt.Errorf("NewPgArtworkRep: %w: %v", ErrPing, err)
			db.Close()
			return
		}
		// Настраиваем пул соединений
		db.SetMaxOpenConns(dbConf.MaxOpenConns)
		db.SetMaxIdleConns(dbConf.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		pgInstance = &PgArtworkRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

func (pg *PgArtworkRep) parseArtworksRows(rows *sql.Rows) ([]*models.Artwork, error) {
	var resArtworks []*models.Artwork
	for rows.Next() {
		var id, authorID, collectionID uuid.UUID
		var title, authorName, collectionTitle, size, material, technic string
		var creationYear, authorBirthYear int
		var authorDeathYear sql.NullInt64
		if err := rows.Scan(&id, &title, &technic, &material, &size, &creationYear,
			&authorID, &authorName, &authorBirthYear, &authorDeathYear,
			&collectionID, &collectionTitle); err != nil {
			return nil, fmt.Errorf("parseArtworksRows: scan error: %v", err)
		}
		deathYear := 0
		if authorDeathYear.Valid {
			deathYear = int(authorDeathYear.Int64)
		}
		author, err := models.NewAuthor(authorID, authorName, authorBirthYear, deathYear)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %v", err)
		}
		collection, err := models.NewCollection(collectionID, collectionTitle)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %v", err)
		}
		user, err := models.NewArtwork(id, title, technic, material, size, creationYear, &author, &collection)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %w: %v", models.ErrValidateArtwork, err)
		}
		resArtworks = append(resArtworks, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("parseArtworksRows: rows iteration error: %v", err)
	}
	return resArtworks, nil
}

func (pg *PgArtworkRep) addFilterParams(query sq.SelectBuilder, filterOps *jsonreqresp.ArtworkFilter) sq.SelectBuilder {
	if filterOps.Title != "" {
		query = query.Where(sq.ILike{"artworks.title": "%" + filterOps.Title + "%"})
	}
	if filterOps.AuthorName != "" {
		query = query.Where(sq.ILike{"author.name": "%" + filterOps.AuthorName + "%"}) // Поиск подстроки
	}
	if filterOps.Collection != "" {
		query = query.Where(sq.ILike{"collection.title": "%" + filterOps.Collection + "%"})
	}
	if filterOps.EventID != uuid.Nil {
		existsSubQuery := sq.Select("1").
			From("Artwork_event ae").
			Where("artworks.id = ae.artworkID").
			Where(sq.Eq{"ae.eventID": filterOps.EventID})
		query = query.Where(sq.Expr("EXISTS (?)", existsSubQuery))
	}
	return query
}

func (pg *PgArtworkRep) addSortParams(query sq.SelectBuilder, sortOps *jsonreqresp.ArtworkSortOps) sq.SelectBuilder {
	switch sortOps.Field {
	case jsonreqresp.TitleSortFieldArtwork:
		query = query.OrderBy("artworks.title " + sortOps.Direction)
	case jsonreqresp.AuthorNameSortFieldArtwork:
		query = query.OrderBy("author.name " + sortOps.Direction)
	case jsonreqresp.CreationYearSortFieldArtwork:
		query = query.OrderBy("artworks.creationYear " + sortOps.Direction)
	}
	return query
}

func (pg *PgArtworkRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Artwork, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseArtworksRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return arts, nil
}

func (pg *PgArtworkRep) GetAllArtworks(ctx context.Context, filterOps *jsonreqresp.ArtworkFilter, sortOps *jsonreqresp.ArtworkSortOps) ([]*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"artworks.id", "artworks.title", "artworks.technic", "artworks.material",
		"artworks.size", "artworks.creationYear",
		"author.id", "author.name", "author.birthyear", "author.deathyear",
		"collection.id", "collection.title").
		From("artworks").
		Join("author ON artworks.authorid = author.id").
		Join("collection ON artworks.collectionid = collection.id")

	query = pg.addFilterParams(query, filterOps)
	query = pg.addSortParams(query, sortOps)
	arts, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgArtworkRep.GetAllArtworks: %w", err)
	}
	return arts, nil
}

func (pg *PgArtworkRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"art.id", "art.title", "art.technic", "art.material",
		"art.size", "art.creationYear",
		"au.id", "au.name", "au.birthyear", "au.deathyear",
		"col.id", "col.title",
	).
		From("artworks art").
		Join("author au ON art.authorid = au.id").
		Join("collection col ON art.collectionid = col.id").
		Where(sq.Eq{"art.id": id})

	arts, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgArtworkRep.GetByID: %w", err)
	}

	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("PgArtworkRep.GetByID: %w", ErrExpectedOneArtwork)
	}
	return arts[0], nil
}

func (pg *PgArtworkRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, querySQL, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: no artowrk added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgArtworkRep) Add(ctx context.Context, e *models.Artwork) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("Artworks").
		Columns("id", "title", "technic", "material", "size", "creationYear", "authorID", "collectionID").
		Values(e.GetID(), e.GetTitle(), e.GetTechnic(), e.GetMaterial(), e.GetSize(), e.GetCreationYear(), e.GetAuthor().GetID(), e.GetCollection().GetID())

	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgArtworkRep.Add: %w", err)
	}
	return nil
}

func (pg *PgArtworkRep) Delete(ctx context.Context, idArt uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("Artworks").
		Where(sq.Eq{"id": idArt})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgArtworkRep.Delete: %w:", err)
	}
	return nil
}

func (pg *PgArtworkRep) Update(ctx context.Context,
	idArt uuid.UUID,
	funcUpdate func(*models.Artwork) (*models.Artwork, error),
) error {
	art, err := pg.GetByID(ctx, idArt)
	if err != nil {
		return fmt.Errorf("PgArtworkRep.Update: %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedArtwork, err := funcUpdate(art)
	if err != nil {
		return fmt.Errorf("PgArtworkRep.Update: %w", ErrUpdateArtwork)
	}
	query := psql.Update("Artworks").
		Set("title", updatedArtwork.GetTitle()).
		Set("material", updatedArtwork.GetMaterial()).
		Set("technic", updatedArtwork.GetTechnic()).
		Set("size", updatedArtwork.GetSize()).
		Set("creationYear", updatedArtwork.GetCreationYear()).
		Set("authorID", updatedArtwork.GetAuthor().GetID()).
		Set("collectionID", updatedArtwork.GetCollection().GetID()).
		Where(sq.Eq{"id": idArt})
	err = pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgArtworkRep.Update %w", err)
	}
	return nil
}

func (pg *PgArtworkRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgArtworkRep) Close() {
	pg.db.Close()
}
