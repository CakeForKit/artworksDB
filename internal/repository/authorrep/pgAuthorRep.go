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
		var birthYear, deathYear int
		if err := rows.Scan(&id, &name, &birthYear, &deathYear); err != nil {
			return nil, fmt.Errorf("parseAuthorsRows: scan error: %v", err)
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

func (pg *PgAuthorRep) GetAllAuthors(ctx context.Context) ([]*models.Author, error) {
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

func (pg *PgAuthorRep) CheckAuthorByID(ctx context.Context, id uuid.UUID) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id").
		From("Author").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("PgAuthorRep.CheckAuthorByID: %w: %v", ErrQueryBuilds, err)
	}
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("PgAuthorRep.CheckAuthorByID: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (pg *PgAuthorRep) AddAuthor(ctx context.Context, e *models.Author) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var deathYear interface{} = e.GetDeathYear()
	if deathYear == 0 {
		deathYear = nil
	}
	query, args, err := psql.Insert("Author").
		Columns("id", "name", "birthYear", "deathYear").
		Values(e.GetID(), e.GetName(), e.GetBirthYear(), deathYear).
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
