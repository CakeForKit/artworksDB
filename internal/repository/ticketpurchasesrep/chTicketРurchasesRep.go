package ticketpurchasesrep

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

type CHTicketPurchasesRep struct {
	db *sql.DB
}

var (
	chInstance *CHTicketPurchasesRep
	chOnce     sync.Once
)

func NewCHTicketPurchasesRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHTicketPurchasesRep, error) {
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
			resErr = fmt.Errorf("%w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHTicketPurchasesRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHTicketPurchasesRep) parseTicketPurchasesRows(rows *sql.Rows) ([]*models.TicketPurchase, error) {
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

func (ch *CHTicketPurchasesRep) GetTPurchasesOfUserID(ctx context.Context, userID uuid.UUID) ([]*models.TicketPurchase, error) {
	query := `
		SELECT tp.id, tp.customerName, tp.customerEmail, 
		       tp.purchaseDate, tp.eventID, tu.userID
		FROM TicketPurchases tp
		JOIN tickets_user tu ON tp.id = tu.ticketID
		WHERE tu.userID = ?`

	rows, err := ch.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("CHTicketPurchasesRep.GetTPurchasesOfUserID: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	res, err := ch.parseTicketPurchasesRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHTicketPurchasesRep.GetTPurchasesOfUserID: %v", err)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (ch *CHTicketPurchasesRep) GetCntTPurchasesForEvent(ctx context.Context, eventID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(tp.id)
		FROM Events e
		LEFT JOIN TicketPurchases tp ON e.id = tp.eventID
		WHERE e.id = ?`

	var count int
	err := ch.db.QueryRowContext(ctx, query, eventID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %w %v", ErrPgTicketPurchasesRep, ErrQueryExec, err)
	}

	return count, nil
}

func (ch *CHTicketPurchasesRep) execChangeQuery(ctx context.Context, query string, args ...interface{}) error {
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

func (ch *CHTicketPurchasesRep) addConnectTicketsUser(ctx context.Context, tp *models.TicketPurchase) error {
	query := "INSERT INTO tickets_user (ticketID, userID) VALUES (?, ?)"
	err := ch.execChangeQuery(ctx, query, tp.GetID(), tp.GetUserID())
	if err != nil {
		return fmt.Errorf("CHTicketPurchasesRep.addConnectTicketsUser: %w", err)
	}
	return nil
}

func (ch *CHTicketPurchasesRep) Add(ctx context.Context, tp *models.TicketPurchase) error {
	query := `
		INSERT INTO TicketPurchases 
		(id, customerName, customerEmail, purchaseDate, eventID) 
		VALUES (?, ?, ?, ?, ?)`

	err := ch.execChangeQuery(ctx, query,
		tp.GetID(),
		tp.GetCustomerName(),
		tp.GetCustomerEmail(),
		tp.GetPurchaseDate(),
		tp.GetEventID(),
	)
	if err != nil {
		return fmt.Errorf("CHTicketPurchasesRep.Add: %w", err)
	}

	if tp.GetUserID() != uuid.Nil {
		if err = ch.addConnectTicketsUser(ctx, tp); err != nil {
			return fmt.Errorf("%w: %v", ErrRowsAffected, err)
		}
	}
	return nil
}

func (ch *CHTicketPurchasesRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHTicketPurchasesRep) Close() {
	ch.db.Close()
}
