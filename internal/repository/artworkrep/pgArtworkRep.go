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
			resErr = fmt.Errorf("%w: %v", ErrOpenConnect, err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			resErr = fmt.Errorf("%w: %v", ErrPing, err)
			db.Close()
			return
		}
		// Настраиваем пул соединений
		db.SetMaxOpenConns(dbConf.MaxIdleConns)
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
		var creationYear, authorBirthYear, authorDeathYear int
		if err := rows.Scan(&id, &title, &technic, &material, &size, &creationYear,
			&authorID, &authorName, &authorBirthYear, &authorDeathYear,
			&collectionID, &collectionTitle); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		author, err := models.NewAuthor(authorID, authorName, authorBirthYear, authorDeathYear)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %v", err)
		}
		collection, err := models.NewCollection(collectionID, collectionTitle)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %v", err)
		}
		user, err := models.NewArtwork(id, title, technic, material, size, creationYear, &author, &collection)
		if err != nil {
			return nil, fmt.Errorf("parseArtworksRows: %v", err)
		}
		resArtworks = append(resArtworks, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return resArtworks, nil
}

func (pg *PgArtworkRep) GetAll(ctx context.Context) ([]*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select(
		"art.id", "art.title", "art.technic", "art.material",
		"art.size", "art.creationYear",
		"au.id", "au.name", "au.birthyear", "au.deathyear",
		"col.id", "col.title",
	).
		From("artworks art").
		Join("author au ON art.authorid = au.id").
		Join("collection col ON art.collectionid = col.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseArtworksRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	}
	return arts, nil
}

func (pg *PgArtworkRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select(
		"art.id", "art.title", "art.technic", "art.material",
		"art.size", "art.creationYear",
		"au.id", "au.name", "au.birthyear", "au.deathyear",
		"col.id", "col.title",
	).
		From("artworks art").
		Join("author au ON art.authorid = au.id").
		Join("collection col ON art.collectionid = col.id").
		Where(sq.Eq{"art.id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseArtworksRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneArtwork, err)
	}
	return arts[0], nil
}

func (pg *PgArtworkRep) GetByTitle(ctx context.Context, title string) ([]*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select(
		"art.id", "art.title", "art.technic", "art.material",
		"art.size", "art.creationYear",
		"au.id", "au.name", "au.birthyear", "au.deathyear",
		"col.id", "col.title",
	).
		From("artworks art").
		Join("author au ON art.authorid = au.id").
		Join("collection col ON art.collectionid = col.id").
		Where(sq.Eq{"art.title": title}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseArtworksRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneArtwork, err)
	}
	return arts, nil
}

func (pg *PgArtworkRep) GetByAuthor(ctx context.Context, author *models.Author) ([]*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select(
		"art.id", "art.title", "art.technic", "art.material",
		"art.size", "art.creationYear",
		"au.id", "au.name", "au.birthyear", "au.deathyear",
		"col.id", "col.title",
	).
		From("artworks art").
		Join("author au ON art.authorid = au.id").
		Join("collection col ON art.collectionid = col.id").
		Where(sq.Eq{"au.id": author.GetID()}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseArtworksRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneArtwork, err)
	}
	return arts, nil
}

func (pg *PgArtworkRep) GetByCreationTime(ctx context.Context, yearBeg int, yearEnd int) ([]*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select(
		"art.id", "art.title", "art.technic", "art.material",
		"art.size", "art.creationYear",
		"au.id", "au.name", "au.birthyear", "au.deathyear",
		"col.id", "col.title",
	).
		From("artworks art").
		Join("author au ON art.authorid = au.id").
		Join("collection col ON art.collectionid = col.id").
		Where(sq.And{
			sq.GtOrEq{"art.creationYear": yearBeg},
			sq.LtOrEq{"art.creationYear": yearEnd},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseArtworksRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneArtwork, err)
	}
	return arts, nil
}

func (pg *PgArtworkRep) GetByEvent(ctx context.Context, event models.Event) ([]*models.Artwork, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	existsSubQuery := sq.Select("1").
		From("Artwork_event ae").
		Where("art.id = ae.artworkID").
		Where(sq.Eq{"ae.eventID": event.GetID()})

	// Основной запрос
	query, args, err := psql.Select(
		"art.id", "art.title", "art.technic", "art.material",
		"art.size", "art.creationYear",
		"au.id", "au.name", "au.birthyear", "au.deathyear",
		"col.id", "col.title",
	).
		From("artworks art").
		Join("author au ON art.authorid = au.id").
		Join("collection col ON art.collectionid = col.id").
		Where(sq.Expr("EXISTS (?)", existsSubQuery)). // Используем Expr для EXISTS
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseArtworksRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrArtworkNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneArtwork, err)
	}
	return arts, nil
}

func (pg *PgArtworkRep) checkAuthorByID(ctx context.Context, id uuid.UUID) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id").
		From("Author").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (pg *PgArtworkRep) addAuthor(ctx context.Context, e *models.Author) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Insert("Author").
		Columns("id", "name", "birthYear", "deathYear").
		Values(e.GetID(), e.GetName(), e.GetBirthYear(), e.GetDeathYear()).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: no Author added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgArtworkRep) checkCollectionByID(ctx context.Context, id uuid.UUID) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id").
		From("Collection").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (pg *PgArtworkRep) addCollection(ctx context.Context, e *models.Collection) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Insert("Collection").
		Columns("id", "title").
		Values(e.GetID(), e.GetTitle()).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: no Collection added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgArtworkRep) Add(ctx context.Context, e *models.Artwork) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	authorExist, err := pg.checkAuthorByID(ctx, e.GetAuthor().GetID())
	if err != nil {
		return fmt.Errorf("check author: %v", err)
	}
	if !authorExist {
		err = pg.addAuthor(ctx, e.GetAuthor())
		if err != nil {
			return fmt.Errorf("add author: %v", err)
		}
	}
	collectionExist, err := pg.checkCollectionByID(ctx, e.GetCollection().GetID())
	if err != nil {
		return fmt.Errorf("check collection: %v", err)
	}
	if !collectionExist {
		err = pg.addCollection(ctx, e.GetCollection())
		if err != nil {
			return fmt.Errorf("add collection: %v", err)
		}
	}
	query, args, err := psql.Insert("Artworks").
		Columns("id", "title", "technic", "material", "size", "creationYear", "authorID", "collectionID").
		Values(e.GetID(), e.GetTitle(), e.GetTechnic(), e.GetMaterial(), e.GetSize(), e.GetCreationYear(), e.GetAuthor().GetID(), e.GetCollection().GetID()).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
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

func (pg *PgArtworkRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Artworks").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: no user with id %s", ErrRowsAffected, id)
	}
	return nil
}

func (pg *PgArtworkRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Artwork) (*models.Artwork, error)) (*models.Artwork, error) {
	user, err := pg.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedArtwork, err := funcUpdate(user)
	if err != nil {
		return nil, fmt.Errorf("funcUpdate: %v", err)
	}
	updatedAuthorID := updatedArtwork.GetAuthor().GetID()
	updatedCollectionID := updatedArtwork.GetCollection().GetID()
	res, err := pg.checkAuthorByID(ctx, updatedAuthorID)
	if err != nil {
		return nil, fmt.Errorf("funcUpdate: %v", err)
	} else if !res {
		return nil, fmt.Errorf("funcUpdate: %v", ErrNoAuthor)
	}
	res, err = pg.checkCollectionByID(ctx, updatedCollectionID)
	if err != nil {
		return nil, fmt.Errorf("funcUpdate: %v", err)
	} else if !res {
		return nil, fmt.Errorf("funcUpdate: %v", ErrNoCollection)
	}
	query, args, err := psql.Update("Artworks").
		Set("title", updatedArtwork.GetTitle()).
		Set("material", updatedArtwork.GetMaterial()).
		Set("technic", updatedArtwork.GetTechnic()).
		Set("size", updatedArtwork.GetSize()).
		Set("creationYear", updatedArtwork.GetCreationYear()).
		Set("authorID", updatedAuthorID).
		Set("collectionID", updatedCollectionID).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("%w: no artworks updated", ErrRowsAffected)
	}
	return updatedArtwork, nil
}

func (pg *PgArtworkRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgArtworkRep) Close() {
	pg.db.Close()
}
