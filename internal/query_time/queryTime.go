package querytime

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type QueryTime interface {
	MeasureTime(start int, stop int, step int, drawGraph bool) error
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
		migrationDir := filepath.Join(projectRoot, "/cmd/measure_time")
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

//	type ExplainResult struct {
//		Plan struct {
//			NodeType           string          `json:"Node Type"`
//			JoinType           string          `json:"Join Type,omitempty"`
//			InnerUnique        bool            `json:"Inner Unique,omitempty"`
//			RelationName       string          `json:"Relation Name,omitempty"`
//			Alias              string          `json:"Alias,omitempty"`
//			RecheckCond        string          `json:"Recheck Cond,omitempty"`
//			ScanDirection      string          `json:"Scan Direction,omitempty"`
//			IndexName          string          `json:"Index Name,omitempty"`
//			IndexCond          string          `json:"Index Cond,omitempty"`
//			ParentRelationship string          `json:"Parent Relationship,omitempty"`
//			Filter             string          `json:"Filter,omitempty"`
//			ActualRows         int             `json:"Actual Rows"`
//			Plans              []ExplainResult `json:"Plans,omitempty"`
//		} `json:"Plan"`
//		PlanningTime  float64 `json:"Planning Time"`
//		ExecutionTime float64 `json:"Execution Time"`
//	}
type PlanNode struct {
	NodeType            string     `json:"Node Type"`
	ParallelAware       bool       `json:"Parallel Aware,omitempty"`
	AsyncCapable        bool       `json:"Async Capable,omitempty"`
	JoinType            string     `json:"Join Type,omitempty"`
	StartupCost         float64    `json:"Startup Cost,omitempty"`
	TotalCost           float64    `json:"Total Cost,omitempty"`
	PlanRows            int        `json:"Plan Rows,omitempty"`
	PlanWidth           int        `json:"Plan Width,omitempty"`
	ActualStartupTime   float64    `json:"Actual Startup Time,omitempty"`
	ActualTotalTime     float64    `json:"Actual Total Time,omitempty"`
	ActualRows          int        `json:"Actual Rows"`
	ActualLoops         int        `json:"Actual Loops,omitempty"`
	InnerUnique         bool       `json:"Inner Unique,omitempty"`
	RelationName        string     `json:"Relation Name,omitempty"`
	Alias               string     `json:"Alias,omitempty"`
	Filter              string     `json:"Filter,omitempty"`
	RowsRemovedByFilter int        `json:"Rows Removed by Filter,omitempty"`
	ScanDirection       string     `json:"Scan Direction,omitempty"`
	IndexName           string     `json:"Index Name,omitempty"`
	IndexCond           string     `json:"Index Cond,omitempty"`
	ParentRelationship  string     `json:"Parent Relationship,omitempty"`
	Plans               []PlanNode `json:"Plans,omitempty"`
}

type ExplainResult struct {
	Plan          PlanNode `json:"Plan"`
	PlanningTime  float64  `json:"Planning Time,omitempty"`
	ExecutionTime float64  `json:"Execution Time,omitempty"`
}

func (q *queryTime) getCountRelationsAE(ctx context.Context) (int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("COUNT(*)").
		From("Artwork_event").ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, fmt.Errorf("no events found")
	}

	var count int
	if err := rows.Scan(&count); err != nil {
		return 0, fmt.Errorf("rows iteration error: %v", err)
	}
	return count, nil
}

func (q *queryTime) getRandomTableID(ctx context.Context, tableName string) (uuid.UUID, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select("id").
		From(tableName).
		OrderBy("random()").
		Limit(1).
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return uuid.Nil, fmt.Errorf("no events found")
	}
	var eventID uuid.UUID
	if err := rows.Scan(&eventID); err != nil {
		return uuid.Nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return eventID, nil
}

func (q *queryTime) addRelationArtworkEvent(ctx context.Context, cntInsert int) error {
	// psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := sq.Expr("SELECT insert_artwork_events($1)", cntInsert).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка построения запроса: %v", err)
	}
	result, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("addRelationArtworkEvent %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("addRelationArtworkEvent %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("addRelationArtworkEvent %w: no Artwork_event added", ErrRowsAffected)
	}

	return nil
}

func parseExplainJSON(resultJSON []byte) ([]ExplainResult, error) {
	// Сначала анмаршалим в промежуточную структуру
	var rawResults []map[string]json.RawMessage
	if err := json.Unmarshal(resultJSON, &rawResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw JSON: %w", err)
	}

	var results []ExplainResult
	for _, rawResult := range rawResults {
		var result ExplainResult

		// Обрабатываем Plan отдельно, так как он рекурсивный
		if planData, ok := rawResult["Plan"]; ok {
			if err := json.Unmarshal(planData, &result.Plan); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Plan: %w", err)
			}
		}

		// Обрабатываем остальные поля
		if planningTime, ok := rawResult["Planning Time"]; ok {
			if err := json.Unmarshal(planningTime, &result.PlanningTime); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Planning Time: %w", err)
			}
		}

		if execTime, ok := rawResult["Execution Time"]; ok {
			if err := json.Unmarshal(execTime, &result.ExecutionTime); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Execution Time: %w", err)
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func printPlan(plan PlanNode, indent int) string {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	res := fmt.Sprintf("%sNode Type: %s\n", indentStr, plan.NodeType)
	// fmt.Printf("%sActual Rows: %d\n", indentStr, plan.ActualRows)

	for _, subPlan := range plan.Plans {
		res += printPlan(subPlan, indent+1)
	}
	return res
}

func (q *queryTime) oneMeasure(ctx context.Context, saveExplain bool, fileExplain *os.File) (*ExplainResult, error) {
	eventID, err := q.getRandomTableID(ctx, "Events")
	if err != nil {
		return nil, fmt.Errorf("oneMeasure: %v", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("Artworks.title").
		From("Artworks").
		Join("Artwork_event ON Artwork_event.artworkID = Artworks.id").
		Where(sq.Eq{"Artwork_event.eventID": eventID}).
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

	results, err := parseExplainJSON(resultJSON)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	resultExplain := results[0]
	// var explainData []ExplainResult
	// if err := json.Unmarshal(resultJSON, &explainData); err != nil {
	// 	return nil, fmt.Errorf("failed to parse JSON: %w", err)
	// }
	// if len(explainData) == 0 {
	// 	return nil, fmt.Errorf("empty explain result")
	// }
	// resultExplain := explainData[0]

	if saveExplain {
		// fmt.Printf("PLAN: %+v\n", resultExplain)
		// if _, err := fileExplain.Write(resultJSON); err != nil {
		// 	return nil, fmt.Errorf("failed to write Explain: %w", err)
		// }
		cnt, err := q.getCountRelationsAE(ctx)
		if err != nil {
			return nil, fmt.Errorf("MeasureTime: %v", err)
		}
		text := fmt.Sprintf("\n%d rows in artworks_event\n", cnt)
		text += printPlan(resultExplain.Plan, 0)
		if _, err := fileExplain.WriteString(text); err != nil {
			return nil, fmt.Errorf("failed to write Explain: %w", err)
		}
	}
	return &resultExplain, nil
}

func (q *queryTime) createIndex() error {
	query := sq.Expr(`
		CREATE INDEX idx_Artwork_event_eventID 
		ON Artwork_event(eventID)
	`)
	sql, _, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("createIndex: %v", err)
	}
	_, err = q.db.Exec(sql)
	if err != nil {
		return fmt.Errorf("createIndex: %v", err)
	}
	return nil
}

func (q *queryTime) dropIndex() error {
	query := sq.Expr(`
		DROP INDEX idx_Artwork_event_eventID;
	`)
	sql, _, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("dropIndex: %v", err)
	}
	_, err = q.db.Exec(sql)
	if err != nil {
		return fmt.Errorf("dropIndex: %v", err)
	}
	return nil
}

func (q *queryTime) MeasureTime(start int, stop int, step int, drawGraph bool) error {
	projectRoot := cnfg.GetProjectRoot()
	dir := filepath.Join(projectRoot, "/measure_results/data/")
	ctx := context.Background()

	cntForOneMeasure := 20
	err := q.addRelationArtworkEvent(ctx, start)
	if err != nil {
		return fmt.Errorf("MeasureTime: %v", err)
	}
	for i := start; i < stop; i += step {
		cnt, err := q.getCountRelationsAE(ctx)
		if err != nil {
			return fmt.Errorf("MeasureTime: %v", err)
		}
		fmt.Printf("Count Artwork_events = %d\n", cnt)

		fnameNotIndex := filepath.Join(dir, fmt.Sprintf("%d_notIndex.txt", i))
		fileNotIndex, err := os.OpenFile(fnameNotIndex, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("MeasureTime: %v", err)
		}
		defer fileNotIndex.Close()
		fnameIndex := filepath.Join(dir, fmt.Sprintf("%d_Index.txt", i))
		fileIndex, err := os.OpenFile(fnameIndex, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("MeasureTime: %v", err)
		}
		defer fileIndex.Close()
		nifileExplain, err := os.OpenFile("./measure_results/notindex_explain.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer nifileExplain.Close()
		ifileExplain, err := os.OpenFile("./measure_results/index_explain.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer ifileExplain.Close()

		for i := range cntForOneMeasure {
			resultExplain, err := q.oneMeasure(ctx, (i == cntForOneMeasure-1), nifileExplain)
			if err != nil {
				return fmt.Errorf("MeasureTime: %v", err)
			}
			if _, err = fileNotIndex.WriteString(fmt.Sprintf("%f\n", resultExplain.ExecutionTime)); err != nil {
				return fmt.Errorf("MeasureTime: %v", err)
			}
			fmt.Printf("ExecutionTime: %f, ActualRows: %d\n", resultExplain.ExecutionTime, resultExplain.Plan.ActualRows)
		}
		err = q.createIndex()
		if err != nil {
			return fmt.Errorf("MeasureTime: %v", err)
		}
		fmt.Print("---------------------------------------\n")
		for i := range cntForOneMeasure {
			resultExplain, err := q.oneMeasure(ctx, (i == cntForOneMeasure-1), ifileExplain)
			if err != nil {
				return fmt.Errorf("MeasureTime: %v", err)
			}
			if _, err = fileIndex.WriteString(fmt.Sprintf("%f\n", resultExplain.ExecutionTime)); err != nil {
				return fmt.Errorf("MeasureTime: %v", err)
			}
			fmt.Printf("ExecutionTime: %f, ActualRows: %d\n", resultExplain.ExecutionTime, resultExplain.Plan.ActualRows)
		}
		err = q.dropIndex()
		if err != nil {
			return fmt.Errorf("MeasureTime: %v", err)
		}

		err = q.addRelationArtworkEvent(ctx, step)
		if err != nil {
			return fmt.Errorf("MeasureTime: %v", err)
		}
	}

	if drawGraph {
		err = DrawGraph(start, stop, step)
		if err != nil {
			return fmt.Errorf("MeasureTime: %v", err)
		}
	}
	return nil
}
