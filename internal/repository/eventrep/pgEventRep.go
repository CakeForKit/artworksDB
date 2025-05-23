package eventrep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
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
		db.SetMaxOpenConns(dbConf.MaxOpenConns)
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
		var canVisit, valid bool
		var cntTickets int
		if err := rows.Scan(&id, &title, &dateBegin, &dateEnd, &canVisit, &address, &cntTickets, &creatorID, &valid); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		user, err := models.NewEvent(id, title, dateBegin, dateEnd, address, canVisit, creatorID, cntTickets, valid, nil)
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

func (pg *PgEventRep) addFilterParams(query sq.SelectBuilder, filterOps *jsonreqresp.EventFilter) sq.SelectBuilder {
	if filterOps.Title != "" {
		query = query.Where(sq.ILike{"events.title": "%" + filterOps.Title + "%"})
	}

	if !filterOps.DateBegin.IsZero() && !filterOps.DateEnd.IsZero() {
		query = query.Where(sq.Or{
			sq.Expr("dateBegin BETWEEN ? AND ?", filterOps.DateBegin, filterOps.DateEnd),
			sq.Expr("dateEnd BETWEEN ? AND ?", filterOps.DateBegin, filterOps.DateEnd),
			sq.And{
				sq.Expr("dateBegin <= ?", filterOps.DateBegin),
				sq.Expr("dateEnd >= ?", filterOps.DateEnd),
			},
		})
	} else if !filterOps.DateBegin.IsZero() {
		query = query.Where(sq.GtOrEq{"dateBegin": filterOps.DateBegin})
	} else if !filterOps.DateEnd.IsZero() {
		query = query.Where(sq.LtOrEq{"dateEnd": filterOps.DateEnd})
	}

	if filterOps.CanVisit != "" {
		canVisit, _ := strconv.ParseBool(filterOps.CanVisit)
		query = query.Where(sq.Eq{"events.canVisit": canVisit})
	}
	return query
}

func (pg *PgEventRep) execQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Event, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := pg.parseEventsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return events, nil
}

func (pg *PgEventRep) GetArtworkIDs(ctx context.Context, eventID uuid.UUID) (uuid.UUIDs, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("artworkID").
		From("Artwork_event").
		Where(sq.Eq{"eventID": eventID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	var artworkIDs uuid.UUIDs
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("PgEventRep.GetArtworkIDs: %v", err)
		}
		artworkIDs = append(artworkIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PgEventRep.GetArtworkIDs rows iteration error: %v", err)
	}
	return artworkIDs, nil
}

func (pg *PgEventRep) joinArtworkIDsToEvents(ctx context.Context, events []*models.Event) ([]*models.Event, error) {
	for _, event := range events {
		artworkIDs, err := pg.GetArtworkIDs(ctx, event.GetID())
		if err != nil {
			return nil, fmt.Errorf("join ArtworkIds %w", err)
		}
		if err := event.AddArtworks(artworkIDs); err != nil {
			return nil, fmt.Errorf("join ArtworkIds %w", err)
		}
	}
	return events, nil
}

func (pg *PgEventRep) GetAll(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"events.id", "events.title", "events.dateBegin", "events.dateEnd", "events.canVisit",
		"events.adress", "events.cntTickets", "events.creatorID", "events.valid").
		From("events")

	query = pg.addFilterParams(query, filterOps)
	events, err := pg.execQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgEventRep.GetAll %w", err)
	}

	events, err = pg.joinArtworkIDsToEvents(ctx, events)
	if err != nil {
		return nil, fmt.Errorf("PgEventRep.GetAll %w", err)
	}

	return events, nil
}

func (pg *PgEventRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"events.id", "events.title", "events.dateBegin", "events.dateEnd", "events.canVisit",
		"events.adress", "events.cntTickets", "events.creatorID", "events.valid").
		From("events").
		Where(sq.Eq{"id": id})

	events, err := pg.execQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgEventRep.GetByID %w", err)
	}
	if len(events) == 0 {
		return nil, nil
	} else if len(events) > 1 {
		return nil, fmt.Errorf("PgEventRep.GetByID %w: %v", ErrExpectedOneEvent, err)
	}
	events, err = pg.joinArtworkIDsToEvents(ctx, events)
	if err != nil {
		return nil, fmt.Errorf("PgEventRep.GetByID %w", err)
	}
	return events[0], nil
}

func (pg *PgEventRep) GetByDate(ctx context.Context, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "title", "dateBegin", "dateEnd", "canVisit", "adress", "cntTickets", "creatorID", "valid").
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
	for _, event := range events {
		artworkIDs, err := pg.GetArtworkIDs(ctx, event.GetID())
		if err != nil {
			return nil, fmt.Errorf("PgEventRep.GetAll %v", err)
		}
		if err := event.AddArtworks(artworkIDs); err != nil {
			return nil, fmt.Errorf("PgEventRep.GetAll %v", err)
		}
	}
	return events, nil
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func (pg *PgEventRep) GetEventsOfArtworkOnDate(ctx context.Context, artworkID uuid.UUID, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	funcCall := sq.DebugSqlizer(sq.Expr("get_event_of_artwork(?, ?, ?)",
		artworkID,
		formatTime(dateBeg),
		formatTime(dateEnd),
	))
	query, args, err := psql.Select("event_id", "title", "dateBegin", "dateEnd", "canVisit", "adress", "cntTickets", "creatorID", "valid").
		From(funcCall).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgEventRep.GetEventsOfArtworkOnDate %w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgEventRep.GetEventsOfArtworkOnDate %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := pg.parseEventsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgEventRep.GetEventsOfArtworkOnDate: %v", err)
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("PgEventRep.GetEventsOfArtworkOnDate: %w", ErrEventNotFound)
	}
	for _, event := range events {
		artworkIDs, err := pg.GetArtworkIDs(ctx, event.GetID())
		if err != nil {
			return nil, fmt.Errorf("PgEventRep.GetAll %v", err)
		}
		if err := event.AddArtworks(artworkIDs); err != nil {
			return nil, fmt.Errorf("PgEventRep.GetAll %v", err)
		}
	}
	return events, nil
}

func (pg *PgEventRep) CheckEmployeeByID(ctx context.Context, id uuid.UUID) (bool, error) {
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
		return fmt.Errorf("%w: no event added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgEventRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Update("Events").
		Set("valid", false).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("PgEventRep.Delete %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgEventRep.Delete %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgEventRep.Delete %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgEventRep.Delete %w: no event deleted", ErrEventNotFound)
	}
	return nil
}

func (pg *PgEventRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Event) (*models.Event, error)) error {
	event, err := pg.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("PgEventRep.Update: %v", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedEvent, err := funcUpdate(event)
	if err != nil {
		return fmt.Errorf("PgEventRep.Update funcUpdate: %v", err)
	}
	query, args, err := psql.Update("Events").
		Set("title", updatedEvent.GetTitle()).
		Set("dateBegin", updatedEvent.GetDateBegin()).
		Set("canVisit", updatedEvent.GetAccess()).
		Set("adress", updatedEvent.GetAddress()).
		Set("cntTickets", updatedEvent.GetTicketCount()).
		Set("creatorID", updatedEvent.GetEmployeeID()).
		Set("valid", updatedEvent.IsValid()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("PgEventRep.Update %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgEventRep.Update %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgEventRep.Update %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgEventRep.Update %w: no event updated", ErrRowsAffected)
	}
	return nil
}

func (pg *PgEventRep) AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkIDs uuid.UUIDs) error {

	for _, artworkID := range artworkIDs {
		psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		query, args, err := psql.Insert("Artwork_event").
			Columns("eventID", "artworkID").
			Values(eventID, artworkID).
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
			return fmt.Errorf("%w: no artwork_event added", ErrRowsAffected)
		}
	}
	return nil
}

func (pg *PgEventRep) DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Artwork_event").
		Where(sq.And{
			sq.Eq{"artworkID": artworkID},
			sq.Eq{"eventID": eventID},
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgEventRep.Delete %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgEventRep.Delete: %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgEventRep.Delete: %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgEventRep.Delete: %w", ErrEventArtowrkNotFound)
	}
	return nil
}

func (pg *PgEventRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgEventRep) Close() {
	pg.db.Close()
}
