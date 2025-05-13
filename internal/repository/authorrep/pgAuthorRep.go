package authorrep

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

type PgAuthorRep struct {
	db *sql.DB
}

var (
	pgInstance *PgAuthorRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect       = errors.New("open connect failed")
	ErrPing              = errors.New("ping failed")
	ErrQueryBuilds       = errors.New("query build failed")
	ErrQueryExec         = errors.New("query execution failed")
	ErrExpectedOneAuthor = errors.New("expected one author")
	ErrRowsAffected      = errors.New("no rows affected")
)

func NewPgAuthorRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (*PgAuthorRep, error) {
	var resErr error
	pgOnce.Do(func() {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			resErr = fmt.Errorf("NewPgAuthorRep: %w: %w", ErrOpenConnect, err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			resErr = fmt.Errorf("NewPgAuthorRep: %w: %w", ErrPing, err)
			db.Close()
			return
		}
		// Настраиваем пул соединений
		db.SetMaxOpenConns(dbConf.MaxIdleConns)
		db.SetMaxIdleConns(dbConf.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		pgInstance = &PgAuthorRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

func (pg *PgAuthorRep) parseAuthorsRows(rows *sql.Rows) ([]*models.Author, error) {
	var resAuthors []*models.Author
	for rows.Next() {
		var id uuid.UUID
		var name string
		var birthYear int
		var authorDeathYear sql.NullInt64
		if err := rows.Scan(&id, &name, &birthYear, &authorDeathYear); err != nil {
			return nil, fmt.Errorf("parseAuthorsRows: scan error: %v", err)
		}
		deathYear := 0
		if authorDeathYear.Valid {
			deathYear = int(authorDeathYear.Int64)
		}
		author, err := models.NewAuthor(id, name, birthYear, deathYear)
		if err != nil {
			return nil, fmt.Errorf("parseAuthorsRows: %v", err)
		}
		resAuthors = append(resAuthors, &author)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return resAuthors, nil
}

func (pg *PgAuthorRep) GetAll(ctx context.Context) ([]*models.Author, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "name", "birthyear", "deathyear").
		From("author").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseAuthorsRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrAuthorNotFound
	}
	return arts, nil
}

func (pg *PgAuthorRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Author, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "name", "birthyear", "deathyear").
		From("Author").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgAuthorRep.GetByID: %w: %v", ErrQueryBuilds, err)
	}
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgAuthorRep.GetByID: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	authors, err := pg.parseAuthorsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgAuthorRep.GetByID %v", err)
	}
	if len(authors) == 0 {
		return nil, ErrAuthorNotFound
	} else if len(authors) > 1 {
		return nil, fmt.Errorf("PgAuthorRep.GetByID %w: %v", ErrAuthorNotFound, err)
	}
	return authors[0], nil
}

func (pg *PgAuthorRep) Add(ctx context.Context, a *models.Author) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var deathYear interface{} = a.GetDeathYear()
	if deathYear == 0 {
		deathYear = nil
	}
	query, args, err := psql.Insert("Author").
		Columns("id", "name", "birthYear", "deathYear").
		Values(a.GetID(), a.GetName(), a.GetBirthYear(), deathYear).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgAuthorRep.CheckAuthorByID: %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgAuthorRep.CheckAuthorByID: %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgAuthorRep.CheckAuthorByID: %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgAuthorRep.CheckAuthorByID: %w: no Author added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgAuthorRep) Delete(ctx context.Context, idAuthor uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Author").
		Where(sq.Eq{"id": idAuthor}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Delete %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Delete %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Delete %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgAuthorRep.Delete %w: no collection with id %s", ErrAuthorNotFound, idAuthor)
	}
	return nil
}

func (pg *PgAuthorRep) Update(
	ctx context.Context,
	idAuthor uuid.UUID,
	funcUpdate func(*models.Author) (*models.Author, error),
) error {
	col, err := pg.GetByID(ctx, idAuthor)
	if err != nil {
		return fmt.Errorf("pgCollectionRep.Update %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedEmployee, err := funcUpdate(col)
	if err != nil {
		return fmt.Errorf("PgEmployeeRep.Update funcUpdate: %v", err)
	}
	query, args, err := psql.Update("Collection").
		Set("name", updatedEmployee.GetName()).
		Set("birthYear", updatedEmployee.GetDeathYear()).
		Set("deathYear", updatedEmployee.GetDeathYear()).
		Where(sq.Eq{"id": idAuthor}).ToSql()
	if err != nil {
		return fmt.Errorf("pgCollectionRep.Update %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("pgCollectionRep.Update %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("pgCollectionRep.Update %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pgCollectionRep.Update %w: no employee added", ErrAuthorNotFound)
	}
	return nil
}

func (pg *PgAuthorRep) HasArtworks(ctx context.Context, authorID uuid.UUID) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("1").
		From("artworks").
		Where(sq.Eq{"authorid": authorID}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("PgAuthorRep.HasArtworks %w: %v", ErrQueryBuilds, err)
	}
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("PgArtworkRep.GetByEvent: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}
