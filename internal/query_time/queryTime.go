package querytime

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type QueryTime interface {
	MeasureTime() error
}

var (
	pgInstance QueryTime
	pgOnceRep  sync.Once
)

var (
	ErrOpenConnect     = errors.New("open connect failed")
	ErrPing            = errors.New("ping failed")
	ErrQueryBuilds     = errors.New("query build failed")
	ErrQueryExec       = errors.New("query execution failed")
	ErrExpectedOneUser = errors.New("expected one user")
	ErrRowsAffected    = errors.New("no rows affected")
)

type queryTime struct {
	db *sql.DB
}

func NewQueryTime() (QueryTime, error) {

	var resErr error
	pgOnceRep.Do(func() {
		// Создение контейнера
		ctx := context.Background()
		_, pgCreds, err := pgtest.GetTestPostgres(ctx)
		if err != nil {
			resErr = fmt.Errorf("queryTime createContaiter: %v", err)
			return
		}
		// -----------
		// Миграции
		projectRoot := cnfg.GetProjectRoot()
		migrationDir := filepath.Join(projectRoot, "migrations")
		err = pgtest.MigrateUp(ctx, migrationDir, &pgCreds)
		if err != nil {
			resErr = fmt.Errorf("queryTime MigrateUp: %v", err)
			return
		}
		// -----------
		// Соединение
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
		// -----------
		// Настраиваем пул соединений
		dbCnfg := cnfg.GetTestDatebaseConfig()
		db.SetMaxOpenConns(dbCnfg.MaxIdleConns)
		db.SetMaxIdleConns(dbCnfg.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbCnfg.ConnMaxLifetime.Hours()))

		pgInstance = &queryTime{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

type ExplainResult struct {
	Plan struct {
		NodeType      string           `json:"Node Type"`
		ParallelAware bool             `json:"Parallel Aware"`
		StartupCost   float64          `json:"Startup Cost"`
		TotalCost     float64          `json:"Total Cost"`
		PlanRows      int              `json:"Plan Rows"`
		PlanWidth     int              `json:"Plan Width"`
		ActualTime    float64          `json:"Actual Total Time"`
		ActualRows    int              `json:"Actual Rows"`
		Plans         []*ExplainResult `json:"Plans,omitempty"`
		// Дополнительные поля
	} `json:"Plan"`
	PlanningTime  float64 `json:"Planning Time"`
	ExecutionTime float64 `json:"Execution Time"`
}

func (q *queryTime) explainQueryJSON() (*ExplainResult, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("a.title", "e.title", "e.dateBegin", "e.dateEnd").
		From("Events e").
		Join("Artwork_event ae ON e.id = ae.eventID").
		Join("artworks a ON ae.artworkID = a.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	query = "EXPLAIN (ANALYZE, FORMAT JSON)  " + query
	var resultJSON []byte
	err = q.db.QueryRowContext(ctx, query, args...).Scan(&resultJSON)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	var explainData []ExplainResult
	if err := json.Unmarshal(resultJSON, &explainData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	if len(explainData) == 0 {
		return nil, fmt.Errorf("empty explain result")
	}

	return &explainData[0], nil
}

func (q *queryTime) MeasureTime() error {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("a.title", "e.title", "e.dateBegin", "e.dateEnd").
		From("Events e").
		Join("Artwork_event ae ON e.id = ae.eventID").
		Join("artworks a ON ae.artworkID = a.id").
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	query = "EXPLAIN (ANALYZE, FORMAT JSON)  " + query
	var resultJSON []byte
	err = q.db.QueryRowContext(ctx, query, args...).Scan(&resultJSON)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	var explainData []ExplainResult
	if err := json.Unmarshal(resultJSON, &explainData); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	if len(explainData) == 0 {
		return fmt.Errorf("empty explain result")
	}
	resultExplain := explainData[0]

	detailedJSON, _ := json.MarshalIndent(resultExplain, "", "  ")
	fmt.Println("\nFull JSON output:")
	fmt.Println(string(detailedJSON))
	return nil
}
