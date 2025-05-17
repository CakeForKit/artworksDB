package eventrep

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

type PgEventRep struct {
	db *sql.DB
}

var (
	pgInstance *PgEventRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect      = errors.New("open connect failed")
	ErrPing             = errors.New("ping failed")
	ErrQueryBuilds      = errors.New("query build failed")
	ErrQueryExec        = errors.New("query execution failed")
	ErrExpectedOneEvent = errors.New("expected one Event")
	ErrRowsAffected     = errors.New("no rows affected")
)

func NewPgEventRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (*PgEventRep, error) {
	var resErr error
	pgOnce.Do(func() {
		// connStr := "postgres://puser:ppassword@postgres_Events:5432/Events"
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

		pgInstance = &PgEventRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

func (pg *PgEventRep) parseEventsRows(rows *sql.Rows) ([]*models.Event, error) {
	var resEvents []*models.Event
	for rows.Next() {
		var id, creatorID uuid.UUID
		var title, address string
		var dateBegin, dateEnd time.Time
		var canVisit bool
		var cntTickets int
		if err := rows.Scan(&id, &title, &dateBegin, &dateEnd, &canVisit, &address, &cntTickets, &creatorID); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		user, err := models.NewEvent(id, title, dateBegin, dateEnd, address, canVisit, creatorID, cntTickets)
		if err != nil {
			return nil, fmt.Errorf("parseEventsRows: %v", err)
		}
		resEvents = append(resEvents, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return resEvents, nil
}

func (pg *PgEventRep) GetAll(ctx context.Context) ([]*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "title", "dateBegin", "dateEnd", "canVisit", "adress", "cntTickets", "creatorID").
		From("events").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := pg.parseEventsRows(rows)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, ErrEventNotFound
	}
	return events, nil
}

func (pg *PgEventRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "title", "dateBegin", "dateEnd", "canVisit", "adress", "cntTickets", "creatorID").
		From("events").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := pg.parseEventsRows(rows)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, ErrEventNotFound
	} else if len(events) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneEvent, err)
	}
	return events[0], nil
}

func (pg *PgEventRep) GetByDate(ctx context.Context, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "title", "dateBegin", "dateEnd", "canVisit", "adress", "cntTickets", "creatorID").
		From("Events").
		Where(sq.Or{
			sq.Expr("dateBegin BETWEEN ? AND ?", dateBeg, dateEnd),
			sq.Expr("dateEnd BETWEEN ? AND ?", dateBeg, dateEnd),
			sq.And{
				sq.Expr("dateBegin <= ?", dateBeg),
				sq.Expr("dateEnd >= ?", dateEnd),
			},
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

	events, err := pg.parseEventsRows(rows)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, ErrEventNotFound
	}
	return events, nil
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func (pg *PgEventRep) GetEventsOfArtworkOnDate(ctx context.Context, artwork *models.Artwork, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	funcCall := sq.DebugSqlizer(sq.Expr("get_event_of_artwork(?, ?, ?)",
		artwork.GetID(),
		formatTime(dateBeg),
		formatTime(dateEnd),
	))
	query, args, err := psql.Select("event_id", "title", "dateBegin", "dateEnd", "canVisit", "adress", "cntTickets", "creatorID").
		From(funcCall).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := pg.parseEventsRows(rows)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, ErrEventNotFound
	}
	return events, nil
}

func (pg *PgEventRep) checkEmployeeByID(ctx context.Context, id uuid.UUID) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id").
		From("Employees").
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

func (pg *PgEventRep) Add(ctx context.Context, e *models.Event) error {
	employeeExist, err := pg.checkEmployeeByID(ctx, e.GetEmployeeID())
	if err != nil {
		return fmt.Errorf("check employee: %v", err)
	} else if !employeeExist {
		return fmt.Errorf("check employee: %v", ErrAddNoEmployee)
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert("Events").
		Columns("id", "title", "dateBegin", "dateEnd", "canVisit", "adress", "cntTickets", "creatorID").
		Values(e.GetID(), e.GetTitle(), e.GetDateBegin(), e.GetDateEnd(), e.GetAccess(), e.GetAddress(), e.GetTicketCount(), e.GetEmployeeID()).
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
		return fmt.Errorf("%w: no employee added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgEventRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Events").
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
		return fmt.Errorf("%w: no event with id %s", ErrRowsAffected, id)
	}
	return nil
}

func (pg *PgEventRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Event) (*models.Event, error)) (*models.Event, error) {
	event, err := pg.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedEvent, err := funcUpdate(event)
	if err != nil {
		return nil, fmt.Errorf("funcUpdate: %v", err)
	}
	query, args, err := psql.Update("Events").
		Set("title", updatedEvent.GetTitle()).
		Set("dateBegin", updatedEvent.GetDateBegin()).
		Set("canVisit", updatedEvent.GetAccess()).
		Set("adress", updatedEvent.GetAddress()).
		Set("cntTickets", updatedEvent.GetTicketCount()).
		Set("creatorID", updatedEvent.GetEmployeeID()).
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
		return nil, fmt.Errorf("%w: no event added", ErrRowsAffected)
	}
	return updatedEvent, nil
}

func (pg *PgEventRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgEventRep) Close() {
	pg.db.Close()
}
