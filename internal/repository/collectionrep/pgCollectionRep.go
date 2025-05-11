package collectionrep

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

type PgCollectionRep struct {
	db *sql.DB
}

var (
	pgInstance *PgCollectionRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect           = errors.New("open connect failed")
	ErrPing                  = errors.New("ping failed")
	ErrQueryBuilds           = errors.New("query build failed")
	ErrQueryExec             = errors.New("query execution failed")
	ErrExpectedOneCollection = errors.New("expected one collection")
	ErrRowsAffected          = errors.New("no rows affected")
)

func NewPgCollectionRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (*PgCollectionRep, error) {
	var resErr error
	pgOnce.Do(func() {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			resErr = fmt.Errorf("NewPgCollectionRep: %w: %w", ErrOpenConnect, err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			resErr = fmt.Errorf("NewPgCollectionRep: %w: %w", ErrPing, err)
			db.Close()
			return
		}
		// Настраиваем пул соединений
		db.SetMaxOpenConns(dbConf.MaxIdleConns)
		db.SetMaxIdleConns(dbConf.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		pgInstance = &PgCollectionRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

func (pg *PgCollectionRep) parseCollectionsRows(rows *sql.Rows) ([]*models.Collection, error) {
	var resCollections []*models.Collection
	for rows.Next() {
		var id uuid.UUID
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			return nil, fmt.Errorf("parseCollectionsRows: scan error: %v", err)
		}
		collection, err := models.NewCollection(id, title)
		if err != nil {
			return nil, fmt.Errorf("parseCollectionsRows: %v", err)
		}
		resCollections = append(resCollections, &collection)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return resCollections, nil
}

func (pg *PgCollectionRep) GetAllCollections(ctx context.Context) ([]*models.Collection, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "title").
		From("collection").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseCollectionsRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, ErrCollectionNotFound
	}
	return arts, nil
}

func (pg *PgCollectionRep) CheckCollectionByID(ctx context.Context, id uuid.UUID) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id").
		From("Collection").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("PgCollectionRep.CheckCollectionByID: %w: %v", ErrQueryBuilds, err)
	}
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("PgCollectionRep.CheckCollectionByID: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (pg *PgCollectionRep) AddCollection(ctx context.Context, e *models.Collection) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Insert("Collection").
		Columns("id", "title").
		Values(e.GetID(), e.GetTitle()).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgCollectionRep.CheckCollectionByID: %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgCollectionRep.CheckCollectionByID: %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgCollectionRep.CheckCollectionByID: %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgCollectionRep.CheckCollectionByID: %w: no Collection added", ErrRowsAffected)
	}
	return nil
}
