package ticketpurchasesrep

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

type PgTicketPurchasesRep struct {
	db *sql.DB
}

var (
	pgInstance *PgTicketPurchasesRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect                = errors.New("open connect failed")
	ErrPing                       = errors.New("ping failed")
	ErrQueryBuilds                = errors.New("query build failed")
	ErrQueryExec                  = errors.New("query execution failed")
	ErrExpectedOneTicketPurchases = errors.New("expected one TicketPurchases")
	ErrRowsAffected               = errors.New("no rows affected")
)

// func NewPgTicketPurchasesRep(ctx context.Context) (TicketPurchasesRep, error) {
func NewPgTicketPurchasesRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (*PgTicketPurchasesRep, error) {
	var resErr error
	pgOnce.Do(func() {
		// connStr := "postgres://puser:ppassword@postgres_artworks:5432/artworks"
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

		pgInstance = &PgTicketPurchasesRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

func (pg *PgTicketPurchasesRep) parseTicketPurchasessRows(rows *sql.Rows) ([]*models.TicketPurchase, error) {
	var resTicketPurchases []*models.TicketPurchase
	for rows.Next() {
		var id, eventID, userID uuid.UUID
		var customerName, customerEmail string
		var purchaseDate time.Time
		if err := rows.Scan(&id, &customerName, &customerEmail, &purchaseDate, &eventID, &userID); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tp, err := models.NewTicketPurchase(id, customerName, customerEmail, purchaseDate, eventID, userID)
		if err != nil {
			return nil, err
		}
		resTicketPurchases = append(resTicketPurchases, &tp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return resTicketPurchases, nil
}

func (pg *PgTicketPurchasesRep) GetTPurchasesOfUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]*models.TicketPurchase, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select(
		"tp.id", "tp.customername", "tp.customeremail",
		"tp.purchasedate", "tp.eventid", "tu.userid",
	).
		From("TicketPurchases tp").
		Join("tickets_user tu ON tp.id = tu.ticketID").
		Where(sq.Eq{"tu.userID": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseTicketPurchasessRows(rows)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return nil, nil
	}
	return arts, nil
}

func (pg *PgTicketPurchasesRep) GetCntTPurchasesForEvent(
	ctx context.Context,
	eventID uuid.UUID,
) (int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select("COUNT(tp.id)").
		From("Events e").
		LeftJoin("TicketPurchases tp ON e.id = tp.eventID").
		Where(sq.Eq{"e.id": eventID}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	var count int
	err = pg.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	return count, nil
}

func (pg *PgTicketPurchasesRep) addConnectTicketsUser(ctx context.Context, tp *models.TicketPurchase) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert("tickets_user").
		Columns("ticketID", "userID").
		Values(tp.GetID(), tp.GetUserID()).
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
		return fmt.Errorf("%w: no tickets_user added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgTicketPurchasesRep) Add(ctx context.Context, tp *models.TicketPurchase) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert("TicketPurchases").
		Columns("id", "customerName", "customerEmail", "purchaseDate", "eventID").
		Values(tp.GetID(), tp.GetCustomerName(), tp.GetCustomerEmail(), tp.GetPurchaseDate(), tp.GetEventID()).
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
		return fmt.Errorf("%w: no TicketPurchases added", ErrRowsAffected)
	}

	if tp.GetUserID() != uuid.Nil {
		if err = pg.addConnectTicketsUser(ctx, tp); err != nil {
			return fmt.Errorf("%w: %v", ErrRowsAffected, err)
		}
	}
	return nil
}

func (pg *PgTicketPurchasesRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgTicketPurchasesRep) Close() {
	pg.db.Close()
}
