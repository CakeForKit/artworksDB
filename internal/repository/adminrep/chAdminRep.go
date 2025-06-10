package adminrep

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

type CHAdminRep struct {
	db *sql.DB
}

var (
	chInstance *CHAdminRep
	chOnce     sync.Once
)

func NewCHAdminRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHAdminRep, error) {
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
			resErr = fmt.Errorf("NewCHAdminRep: %w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHAdminRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHAdminRep) parseAdminsRows(rows *sql.Rows) ([]*models.Admin, error) {
	var resAdmins []*models.Admin
	for rows.Next() {
		var id uuid.UUID
		var username, login, hashedPassword string
		var createdAt time.Time
		var valid uint8
		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &valid); err != nil {
			return nil, fmt.Errorf("parseAdminsRows, scan error: %v", err)
		}
		admin, err := models.NewAdmin(id, username, login, hashedPassword, createdAt, valid == 1)
		if err != nil {
			return nil, err
		}
		resAdmins = append(resAdmins, &admin)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("parseAdminsRows, rows iteration error: %v", err)
	}
	return resAdmins, nil
}

func (ch *CHAdminRep) execSelectQuery(ctx context.Context, query string, args ...interface{}) ([]*models.Admin, error) {
	rows, err := ch.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	res, err := ch.parseAdminsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

func (ch *CHAdminRep) GetAll(ctx context.Context) ([]*models.Admin, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, valid FROM Admins"
	res, err := ch.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CHAdminRep.GetAll: %v", err)
	}
	return res, nil
}

func (ch *CHAdminRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Admin, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, valid FROM Admins WHERE id = ?"
	res, err := ch.execSelectQuery(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("CHAdminRep.GetByID: %v", err)
	}
	if len(res) == 0 {
		return nil, ErrAdminNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("CHAdminRep.GetByID: %w: %v", ErrExpectedOneAdmin, err)
	}
	return res[0], nil
}

func (ch *CHAdminRep) GetByLogin(ctx context.Context, login string) (*models.Admin, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, valid FROM Admins WHERE login = ?"
	res, err := ch.execSelectQuery(ctx, query, login)
	if err != nil {
		return nil, fmt.Errorf("CHAdminRep.GetByLogin: %v", err)
	}
	if len(res) == 0 {
		return nil, ErrAdminNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("CHAdminRep.GetByLogin: %w: %v", ErrExpectedOneAdmin, err)
	}
	return res[0], nil
}

func (ch *CHAdminRep) execChangeQuery(ctx context.Context, query string, args ...interface{}) error {
	result, err := ch.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}

	// Check affected rows (ClickHouse doesn't support RowsAffected natively)
	// This is a workaround for ClickHouse
	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return fmt.Errorf("%w: no rows modified", ErrRowsAffected)
	}
	return nil
}

func (ch *CHAdminRep) Add(ctx context.Context, e *models.Admin) error {
	_, err := ch.GetByLogin(ctx, e.GetLogin())
	if err == nil {
		return ErrDuplicateLoginAdm
	} else if err != ErrAdminNotFound {
		return fmt.Errorf("CHAdminRep.Add %v", err)
	}

	query := `INSERT INTO Admins (id, username, login, hashedPassword, createdAt, valid) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	valid := uint8(0)
	if e.IsValid() {
		valid = 1
	}

	err = ch.execChangeQuery(ctx, query,
		e.GetID(),
		e.GetUsername(),
		e.GetLogin(),
		e.GetHashedPassword(),
		e.GetCreatedAt(),
		valid)

	if err != nil {
		return fmt.Errorf("CHAdminRep.Add: %w", err)
	}
	return nil
}

func (ch *CHAdminRep) Delete(ctx context.Context, id uuid.UUID) error {
	query := "ALTER TABLE Admins DELETE WHERE id = ?"
	err := ch.execChangeQuery(ctx, query, id)
	if err != nil {
		return fmt.Errorf("CHAdminRep.Delete: %w", err)
	}
	return nil
}

func (ch *CHAdminRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Admin) (*models.Admin, error),
) error {
	admin, err := ch.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("CHAdminRep.Update %w", err)
	}

	updatedAdmin, err := funcUpdate(admin)
	if err != nil {
		return fmt.Errorf("CHAdminRep.Update: %w", ErrUpdateAdmin)
	}

	valid := uint8(0)
	if updatedAdmin.IsValid() {
		valid = 1
	}

	query := `ALTER TABLE Admins UPDATE 
		username = ?, 
		login = ?, 
		hashedPassword = ?, 
		valid = ? 
		WHERE id = ?`

	err = ch.execChangeQuery(ctx, query,
		updatedAdmin.GetUsername(),
		updatedAdmin.GetLogin(),
		updatedAdmin.GetHashedPassword(),
		valid,
		id)

	if err != nil {
		return fmt.Errorf("CHAdminRep.Update: %w", err)
	}
	return nil
}

func (ch *CHAdminRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHAdminRep) Close() {
	ch.db.Close()
}
