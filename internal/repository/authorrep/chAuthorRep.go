package authorrep

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

type CHAuthorRep struct {
	db *sql.DB
}

var (
	chInstance *CHAuthorRep
	chOnce     sync.Once
)

func NewCHAuthorRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHAuthorRep, error) {
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
			resErr = fmt.Errorf("NewCHAuthorRep: %w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHAuthorRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHAuthorRep) parseAuthorsRows(rows *sql.Rows) ([]*models.Author, error) {
	var resAuthors []*models.Author
	for rows.Next() {
		var id uuid.UUID
		var name string
		var birthYear int32
		var deathYear sql.NullInt32
		if err := rows.Scan(&id, &name, &birthYear, &deathYear); err != nil {
			return nil, fmt.Errorf("parseAuthorsRows: scan error: %v", err)
		}

		deathYearValue := 0
		if deathYear.Valid {
			deathYearValue = int(deathYear.Int32)
		}

		author, err := models.NewAuthor(id, name, int(birthYear), deathYearValue)
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

func (ch *CHAuthorRep) execSelectQuery(ctx context.Context, query string, args ...interface{}) ([]*models.Author, error) {
	rows, err := ch.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	res, err := ch.parseAuthorsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

func (ch *CHAuthorRep) GetAll(ctx context.Context) ([]*models.Author, error) {
	query := "SELECT id, name, birthYear, deathYear FROM Author"
	res, err := ch.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CHAuthorRep.GetAll: %v", err)
	}
	return res, nil
}

func (ch *CHAuthorRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Author, error) {
	query := "SELECT id, name, birthYear, deathYear FROM Author WHERE id = ?"
	res, err := ch.execSelectQuery(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("CHAuthorRep.GetByID: %v", err)
	}

	if len(res) == 0 {
		return nil, ErrAuthorNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("CHAuthorRep.GetByID %w", ErrExpectedOneAuthor)
	}
	return res[0], nil
}

func (ch *CHAuthorRep) execChangeQuery(ctx context.Context, query string, args ...interface{}) error {
	result, err := ch.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	// ClickHouse doesn't fully support RowsAffected, but we can still check for errors
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHAuthorRep) Add(ctx context.Context, a *models.Author) error {
	query := "INSERT INTO Author (id, name, birthYear, deathYear) VALUES (?, ?, ?, ?)"

	var deathYear interface{} = nil
	if a.GetDeathYear() != 0 {
		deathYear = a.GetDeathYear()
	}

	err := ch.execChangeQuery(ctx, query,
		a.GetID(),
		a.GetName(),
		a.GetBirthYear(),
		deathYear)

	if err != nil {
		return fmt.Errorf("CHAuthorRep.Add: %w", err)
	}
	return nil
}

func (ch *CHAuthorRep) Delete(ctx context.Context, idAuthor uuid.UUID) error {
	query := "ALTER TABLE Author DELETE WHERE id = ?"
	err := ch.execChangeQuery(ctx, query, idAuthor)
	if err != nil {
		return fmt.Errorf("CHAuthorRep.Delete: %w", err)
	}
	return nil
}

func (ch *CHAuthorRep) Update(
	ctx context.Context,
	idAuthor uuid.UUID,
	funcUpdate func(*models.Author) (*models.Author, error),
) error {
	author, err := ch.GetByID(ctx, idAuthor)
	if err != nil {
		return fmt.Errorf("CHAuthorRep.Update %w", err)
	}

	updatedAuthor, err := funcUpdate(author)
	if err != nil {
		return fmt.Errorf("CHAuthorRep.Update: %w", ErrUpdateAuthor)
	}

	query := "ALTER TABLE Author UPDATE name = ?, birthYear = ?, deathYear = ? WHERE id = ?"

	var deathYear interface{} = nil
	if updatedAuthor.GetDeathYear() != 0 {
		deathYear = updatedAuthor.GetDeathYear()
	}

	err = ch.execChangeQuery(ctx, query,
		updatedAuthor.GetName(),
		updatedAuthor.GetBirthYear(),
		deathYear,
		idAuthor)

	if err != nil {
		return fmt.Errorf("CHAuthorRep.Update: %w", err)
	}
	return nil
}

func (ch *CHAuthorRep) HasArtworks(ctx context.Context, authorID uuid.UUID) (bool, error) {
	query := "SELECT 1 FROM artworks WHERE authorID = ? LIMIT 1"
	rows, err := ch.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return false, fmt.Errorf("CHAuthorRep.HasArtworks: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (ch *CHAuthorRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHAuthorRep) Close() {
	ch.db.Close()
}
