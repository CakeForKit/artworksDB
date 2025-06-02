package collectionrep

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
)

type CHCollectionRep struct {
	db *sql.DB
}

var (
	chInstance *CHCollectionRep
	chOnce     sync.Once
)

func NewCHCollectionRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHCollectionRep, error) {
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
			resErr = fmt.Errorf("NewCHCollectionRep: %w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHCollectionRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHCollectionRep) parseCollectionsRows(rows *sql.Rows) ([]*models.Collection, error) {
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

func (ch *CHCollectionRep) execSelectQuery(ctx context.Context, query string, args ...interface{}) ([]*models.Collection, error) {
	rows, err := ch.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	res, err := ch.parseCollectionsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

func (ch *CHCollectionRep) GetAllCollections(ctx context.Context) ([]*models.Collection, error) {
	query := "SELECT id, title FROM Collection"
	res, err := ch.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CHCollectionRep.GetAllCollections: %v", err)
	}
	return res, nil
}

func (ch *CHCollectionRep) GetCollectionByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	query := "SELECT id, title FROM Collection WHERE id = ?"
	res, err := ch.execSelectQuery(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("CHCollectionRep.GetCollectionByID: %v", err)
	}
	if len(res) == 0 {
		return nil, ErrCollectionNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("CHCollectionRep.GetCollectionByID: %w", ErrExpectedOneCollection)
	}
	return res[0], nil
}

func (ch *CHCollectionRep) execChangeQuery(ctx context.Context, query string, args ...interface{}) error {
	result, err := ch.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support, but we can still check for errors
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHCollectionRep) AddCollection(ctx context.Context, c *models.Collection) error {
	query := "INSERT INTO Collection (id, title) VALUES (?, ?)"

	err := ch.execChangeQuery(ctx, query,
		c.GetID(),
		c.GetTitle())

	if err != nil {
		return fmt.Errorf("CHCollectionRep.AddCollection: %w", err)
	}
	return nil
}

func (ch *CHCollectionRep) DeleteCollection(ctx context.Context, idCol uuid.UUID) error {
	query := "ALTER TABLE Collection DELETE WHERE id = ?"
	err := ch.execChangeQuery(ctx, query, idCol)
	if err != nil {
		return fmt.Errorf("CHCollectionRep.DeleteCollection: %w", err)
	}
	return nil
}

func (ch *CHCollectionRep) UpdateCollection(
	ctx context.Context,
	idCol uuid.UUID,
	funcUpdate func(*models.Collection) (*models.Collection, error),
) error {
	col, err := ch.GetCollectionByID(ctx, idCol)
	if err != nil {
		return fmt.Errorf("CHCollectionRep.UpdateCollection %w", err)
	}

	updatedCollection, err := funcUpdate(col)
	if err != nil {
		return fmt.Errorf("CHCollectionRep.UpdateCollection: %w", ErrUpdateCollection)
	}

	query := "ALTER TABLE Collection UPDATE title = ? WHERE id = ?"
	err = ch.execChangeQuery(ctx, query,
		updatedCollection.GetTitle(),
		idCol)

	if err != nil {
		return fmt.Errorf("CHCollectionRep.UpdateCollection: %w", err)
	}
	return nil
}

func (ch *CHCollectionRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHCollectionRep) Close() {
	ch.db.Close()
}
