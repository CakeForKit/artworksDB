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
		db.SetMaxOpenConns(dbConf.MaxOpenConns)
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

func (pg *PgAuthorRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Author, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	res, err := pg.parseAuthorsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

func (pg *PgAuthorRep) GetAll(ctx context.Context) ([]*models.Author, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("id", "name", "birthyear", "deathyear").
		From("author")
	res, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgAuthorRep.GetAll: %v", err)
	}
	return res, nil
}

func (pg *PgAuthorRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Author, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("id", "name", "birthyear", "deathyear").
		From("Author").
		Where(sq.Eq{"id": id})
	res, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgAuthorRep.GetByID: %v", err)
	}

	if len(res) == 0 {
		return nil, ErrAuthorNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("PgAuthorRep.GetByID %w", ErrExpectedOneAuthor)
	}
	return res[0], nil
}

func (pg *PgAuthorRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
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
		return fmt.Errorf("%w: no added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgAuthorRep) Add(ctx context.Context, a *models.Author) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var deathYear interface{} = a.GetDeathYear()
	if deathYear == 0 {
		deathYear = nil
	}
	query := psql.Insert("Author").
		Columns("id", "name", "birthYear", "deathYear").
		Values(a.GetID(), a.GetName(), a.GetBirthYear(), deathYear)
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Add: %w", err)
	}
	return nil
}

func (pg *PgAuthorRep) Delete(ctx context.Context, idAuthor uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("Author").
		Where(sq.Eq{"id": idAuthor})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgAuthorRep) Update(
	ctx context.Context,
	idAuthor uuid.UUID,
	funcUpdate func(*models.Author) (*models.Author, error),
) error {
	author, err := pg.GetByID(ctx, idAuthor)
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Update %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedAuthor, err := funcUpdate(author)
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Update: %w", ErrUpdateAuthor)
	}

	query := psql.Update("Author").
		Set("name", updatedAuthor.GetName()).
		Set("birthYear", updatedAuthor.GetBirthYear())
		// Устанавливаем deathYear в NULL если значение равно 0
	if updatedAuthor.GetDeathYear() == 0 {
		query = query.Set("deathYear", nil)
	} else {
		query = query.Set("deathYear", updatedAuthor.GetDeathYear())
	}
	query = query.Where(sq.Eq{"id": idAuthor})
	err = pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgAuthorRep.Update: %w", err)
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
		return false, fmt.Errorf("PgAuthorRep.HasArtworks: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}
