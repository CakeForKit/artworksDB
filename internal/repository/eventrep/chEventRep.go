package eventrep

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
)

type CHEventRep struct {
	db *sql.DB
}

var (
	chInstance *CHEventRep
	chOnce     sync.Once
)

func NewCHEventRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHEventRep, error) {
	// fmt.Printf("NewCHEventRep\n\n")
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
			resErr = fmt.Errorf("NewCHEventRep: %w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHEventRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHEventRep) parseEventsRows(rows *sql.Rows) ([]*models.Event, error) {
	var resEvents []*models.Event
	for rows.Next() {
		var id, creatorID uuid.UUID
		var title, address string
		var dateBegin, dateEnd time.Time
		var canVisit, valid uint8
		var cntTickets int32

		if err := rows.Scan(&id, &title, &dateBegin, &dateEnd, &canVisit, &address, &cntTickets, &creatorID, &valid); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		event, err := models.NewEvent(
			id,
			title,
			dateBegin,
			dateEnd,
			address,
			canVisit == 1,
			creatorID,
			int(cntTickets),
			valid == 1,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("parseEventsRows: %v", err)
		}
		resEvents = append(resEvents, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return resEvents, nil
}

func (ch *CHEventRep) buildFilterConditions(filterOps *jsonreqresp.EventFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filterOps.Title != "" {
		conditions = append(conditions, "Events.title LIKE ?")
		args = append(args, "%"+filterOps.Title+"%")
	}

	if !filterOps.DateBegin.IsZero() && !filterOps.DateEnd.IsZero() {
		conditions = append(conditions,
			"(dateBegin BETWEEN ? AND ? OR dateEnd BETWEEN ? AND ? OR (dateBegin <= ? AND dateEnd >= ?))")
		args = append(args,
			filterOps.DateBegin, filterOps.DateEnd,
			filterOps.DateBegin, filterOps.DateEnd,
			filterOps.DateBegin, filterOps.DateEnd,
		)
	} else if !filterOps.DateBegin.IsZero() {
		conditions = append(conditions, "dateBegin >= ?")
		args = append(args, filterOps.DateBegin)
	} else if !filterOps.DateEnd.IsZero() {
		conditions = append(conditions, "dateEnd <= ?")
		args = append(args, filterOps.DateEnd)
	}

	if filterOps.CanVisit != "" {
		canVisit := uint8(0)
		if filterOps.CanVisit == "true" {
			canVisit = 1
		}
		conditions = append(conditions, "Events.canVisit = ?")
		args = append(args, canVisit)
	}

	conditions = append(conditions, "Events.valid = 1")

	if len(conditions) == 0 {
		return "", nil
	}

	return "WHERE " + joinConditions(conditions, " AND "), args
}

func joinConditions(conditions []string, sep string) string {
	if len(conditions) == 0 {
		return ""
	}
	if len(conditions) == 1 {
		return conditions[0]
	}

	result := conditions[0]
	for _, cond := range conditions[1:] {
		result += sep + cond
	}
	return result
}

func (ch *CHEventRep) GetArtworkIDs(ctx context.Context, eventID uuid.UUID) (uuid.UUIDs, error) {
	query := "SELECT artworkID FROM Artwork_event WHERE eventID = ?"
	rows, err := ch.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetArtworkIDs: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	var artworkIDs uuid.UUIDs
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("CHEventRep.GetArtworkIDs: %v", err)
		}
		artworkIDs = append(artworkIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CHEventRep.GetArtworkIDs rows iteration error: %v", err)
	}
	return artworkIDs, nil
}

func (ch *CHEventRep) joinArtworkIDsToEvents(ctx context.Context, events []*models.Event) ([]*models.Event, error) {
	for _, event := range events {
		artworkIDs, err := ch.GetArtworkIDs(ctx, event.GetID())
		if err != nil {
			return nil, fmt.Errorf("join ArtworkIds %w", err)
		}
		if err := event.AddArtworks(artworkIDs); err != nil {
			return nil, fmt.Errorf("join ArtworkIds %w", err)
		}
	}
	return events, nil
}

func (ch *CHEventRep) GetAll(ctx context.Context, filterOps *jsonreqresp.EventFilter) ([]*models.Event, error) {
	baseQuery := `
		SELECT 
			id, title, dateBegin, dateEnd, canVisit, 
			adress, cntTickets, creatorID, valid
		FROM Events`

	filterClause, filterArgs := ch.buildFilterConditions(filterOps)

	query := baseQuery
	if filterClause != "" {
		query += " " + filterClause
	}

	rows, err := ch.db.QueryContext(ctx, query, filterArgs...)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetAll %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := ch.parseEventsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetAll %w", err)
	}

	events, err = ch.joinArtworkIDsToEvents(ctx, events)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetAll %w", err)
	}

	return events, nil
}

func (ch *CHEventRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	query := `
		SELECT 
			id, title, dateBegin, dateEnd, canVisit, 
			adress, cntTickets, creatorID, valid
		FROM Events
		WHERE id = ?`

	rows, err := ch.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetByID %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := ch.parseEventsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetByID %w", err)
	}

	if len(events) == 0 {
		return nil, ErrEventNotFound
	} else if len(events) > 1 {
		return nil, fmt.Errorf("CHEventRep.GetByID %w", ErrExpectedOneEvent)
	}

	events, err = ch.joinArtworkIDsToEvents(ctx, events)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetByID %w", err)
	}

	return events[0], nil
}

func (ch *CHEventRep) GetEventsOfArtworkOnDate(ctx context.Context, artworkID uuid.UUID, dateBeg time.Time, dateEnd time.Time) ([]*models.Event, error) {
	query := `
		SELECT 
			e.id, e.title, e.dateBegin, e.dateEnd, e.canVisit, 
			e.adress, e.cntTickets, e.creatorID, e.valid
		FROM Events e
		JOIN Artwork_event ae ON e.id = ae.eventID
		WHERE ae.artworkID = ?
		AND e.dateBegin <= ?
		AND e.dateEnd >= ?`

	rows, err := ch.db.QueryContext(ctx, query, artworkID, dateEnd, dateBeg)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetEventsOfArtworkOnDate %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	events, err := ch.parseEventsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetEventsOfArtworkOnDate: %v", err)
	}

	if len(events) == 0 {
		return nil, ErrEventNotFound
	}

	for _, event := range events {
		artworkIDs, err := ch.GetArtworkIDs(ctx, event.GetID())
		if err != nil {
			return nil, fmt.Errorf("CHEventRep.GetEventsOfArtworkOnDate %v", err)
		}
		if err := event.AddArtworks(artworkIDs); err != nil {
			return nil, fmt.Errorf("CHEventRep.GetEventsOfArtworkOnDate %v", err)
		}
	}

	return events, nil
}

func (ch *CHEventRep) GetCollectionsStat(ctx context.Context, eventID uuid.UUID) ([]*models.StatCollections, error) {
	query := `
		SELECT 
			c.id AS collection_id,
			c.title AS collection_title,
			COUNT(a.id) AS artwork_count
		FROM Artwork_event ae
		JOIN Artworks a ON ae.artworkID = a.id
		JOIN Collection c ON a.collectionID = c.id
		WHERE ae.eventID = ?
		GROUP BY c.id, c.title
		ORDER BY artwork_count DESC`

	rows, err := ch.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("CHEventRep.GetCollectionsStat %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	var res []*models.StatCollections
	for rows.Next() {
		var colID uuid.UUID
		var colTitle string
		var cntArtworks int64

		if err := rows.Scan(&colID, &colTitle, &cntArtworks); err != nil {
			return nil, fmt.Errorf("CHEventRep.GetCollectionsStat scan error: %v", err)
		}

		statCol, err := models.NewStatCollections(colID, colTitle, int(cntArtworks))
		if err != nil {
			return nil, fmt.Errorf("CHEventRep.GetCollectionsStat: %w", err)
		}
		res = append(res, &statCol)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CHEventRep.GetCollectionsStat rows iteration error: %v", err)
	}

	return res, nil
}

func (ch *CHEventRep) CheckEmployeeByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := "SELECT id FROM Employees WHERE id = ? LIMIT 1"
	rows, err := ch.db.QueryContext(ctx, query, id)
	if err != nil {
		return false, fmt.Errorf("CHEventRep.CheckEmployeeByID: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (ch *CHEventRep) execChangeQuery(ctx context.Context, query string, args ...interface{}) error {
	result, err := ch.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHEventRep) Add(ctx context.Context, e *models.Event) error {
	query := `
		INSERT INTO Events 
		(id, title, dateBegin, dateEnd, canVisit, adress, cntTickets, creatorID, valid) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)`

	canVisit := uint8(0)
	if e.GetAccess() {
		canVisit = 1
	}

	err := ch.execChangeQuery(ctx, query,
		e.GetID(),
		e.GetTitle(),
		e.GetDateBegin(),
		e.GetDateEnd(),
		canVisit,
		e.GetAddress(),
		e.GetTicketCount(),
		e.GetEmployeeID(),
	)
	if err != nil {
		return fmt.Errorf("CHEventRep.Add: %w", err)
	}
	return nil
}

func (ch *CHEventRep) Delete(ctx context.Context, id uuid.UUID) error {
	query := "ALTER TABLE Events UPDATE valid = 0 WHERE id = ?"
	err := ch.execChangeQuery(ctx, query, id)
	if err != nil {
		return fmt.Errorf("CHEventRep.Delete: %w", err)
	}
	return nil
}

func (ch *CHEventRep) RealDelete(ctx context.Context, id uuid.UUID) error {
	query := "ALTER TABLE Events DELETE WHERE id = ?"
	err := ch.execChangeQuery(ctx, query, id)
	if err != nil {
		return fmt.Errorf("CHEventRep.RealDelete: %w", err)
	}
	return nil
}

func (ch *CHEventRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Event) (*models.Event, error)) error {
	event, err := ch.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("CHEventRep.Update: %v", err)
	}

	updatedEvent, err := funcUpdate(event)
	if err != nil {
		return fmt.Errorf("CHEventRep.Update: %w", ErrUpdateEvent)
	}

	canVisit := uint8(0)
	if updatedEvent.GetAccess() {
		canVisit = 1
	}

	query := `
		ALTER TABLE Events UPDATE 
		title = ?, 
		dateBegin = ?, 
		dateEnd = ?, 
		canVisit = ?, 
		adress = ?, 
		cntTickets = ?, 
		creatorID = ? 
		WHERE id = ?`

	err = ch.execChangeQuery(ctx, query,
		updatedEvent.GetTitle(),
		updatedEvent.GetDateBegin(),
		updatedEvent.GetDateEnd(),
		canVisit,
		updatedEvent.GetAddress(),
		updatedEvent.GetTicketCount(),
		updatedEvent.GetEmployeeID(),
		id,
	)
	if err != nil {
		return fmt.Errorf("CHEventRep.Update %w", err)
	}
	return nil
}

func (ch *CHEventRep) AddArtworksToEvent(ctx context.Context, eventID uuid.UUID, artworkIDs uuid.UUIDs) error {
	for _, artworkID := range artworkIDs {
		query := "INSERT INTO Artwork_event (eventID, artworkID) VALUES (?, ?)"
		err := ch.execChangeQuery(ctx, query, eventID, artworkID)
		if err != nil {
			return fmt.Errorf("CHEventRep.AddArtworksToEvent %w", ErrEventArtowrkNotFound)
		}
	}
	return nil
}

func (ch *CHEventRep) DeleteArtworkFromEvent(ctx context.Context, eventID uuid.UUID, artworkID uuid.UUID) error {
	query := "ALTER TABLE Artwork_event DELETE WHERE eventID = ? AND artworkID = ?"
	err := ch.execChangeQuery(ctx, query, eventID, artworkID)
	if err != nil {
		return fmt.Errorf("CHEventRep.DeleteArtworkFromEvent %w", ErrEventArtowrkNotFound)
	}
	return nil
}

func (ch *CHEventRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHEventRep) Close() {
	ch.db.Close()
}
