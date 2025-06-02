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

func NewPgCollectionRep(ctx context.Context, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (*PgCollectionRep, error) {
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
		db.SetMaxOpenConns(dbConf.MaxOpenConns)
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

func (pg *PgCollectionRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Collection, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	res, err := pg.parseCollectionsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

func (pg *PgCollectionRep) GetAllCollections(ctx context.Context) ([]*models.Collection, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("id", "title").
		From("collection")
	res, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgCollectionRep.GetAll: %v", err)
	}
	return res, nil
}

func (pg *PgCollectionRep) GetCollectionByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("id", "title").
		From("Collection").
		Where(sq.Eq{"id": id})
	res, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgCollectionRep.GetByID: %v", err)
	}
	if len(res) == 0 {
		return nil, ErrCollectionNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("PgCollectionRep.GetByID: %w", ErrExpectedOneCollection)
	}
	return res[0], nil
}

func (pg *PgCollectionRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
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

func (pg *PgCollectionRep) AddCollection(ctx context.Context, e *models.Collection) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Insert("Collection").
		Columns("id", "title").
		Values(e.GetID(), e.GetTitle())
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgCollectionRep.Add: %w", err)
	}
	return nil
}

func (pg *PgCollectionRep) DeleteCollection(ctx context.Context, idCol uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("Collection").
		Where(sq.Eq{"id": idCol})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgCollectionRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgCollectionRep) UpdateCollection(
	ctx context.Context,
	idCol uuid.UUID,
	funcUpdate func(*models.Collection) (*models.Collection, error),
) error {
	col, err := pg.GetCollectionByID(ctx, idCol)
	if err != nil {
		return fmt.Errorf("pgCollectionRep.Update %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedEmployee, err := funcUpdate(col)
	if err != nil {
		return fmt.Errorf("pgCollectionRep.Update: %w", ErrUpdateCollection)
	}
	query := psql.Update("Collection").
		Set("title", updatedEmployee.GetTitle()).
		Where(sq.Eq{"id": idCol})
	err = pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("pgCollectionRep.Update: %w", err)
	}
	return nil
}
