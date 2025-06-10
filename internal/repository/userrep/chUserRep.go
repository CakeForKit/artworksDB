package userrep

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

type CHUserRep struct {
	db *sql.DB
}

var (
	chInstance *CHUserRep
	chOnce     sync.Once
)

func NewCHUserRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHUserRep, error) {
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
			resErr = fmt.Errorf("NewCHUserRep: %w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHUserRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHUserRep) parseUsersRows(rows *sql.Rows) ([]*models.User, error) {
	var resUsers []*models.User
	for rows.Next() {
		var id uuid.UUID
		var username, login, hashedPassword string
		var email sql.NullString
		var createdAt time.Time
		var subscribeMail uint8

		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &email, &subscribeMail); err != nil {
			return nil, fmt.Errorf("parseUsersRows scan error: %v", err)
		}

		emailValue := ""
		if email.Valid {
			emailValue = email.String
		}

		user, err := models.NewUser(id, username, login, hashedPassword, createdAt, emailValue, subscribeMail == 1)
		if err != nil {
			return nil, fmt.Errorf("parseUsersRows: %v", err)
		}
		resUsers = append(resUsers, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("parseUsersRows rows iteration error: %v", err)
	}
	return resUsers, nil
}

func (ch *CHUserRep) GetAll(ctx context.Context) ([]*models.User, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, email, subscribeMail FROM Users"
	rows, err := ch.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetAll: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	users, err := ch.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetAll: %v", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users, nil
}

func (ch *CHUserRep) GetAllSubscribed(ctx context.Context) ([]*models.User, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, email, subscribeMail FROM Users WHERE subscribeMail = 1"
	rows, err := ch.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetAllSubscribed: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	users, err := ch.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetAllSubscribed: %v", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users, nil
}

func (ch *CHUserRep) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, email, subscribeMail FROM Users WHERE id = ?"
	rows, err := ch.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetByID: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	users, err := ch.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetByID: %v", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	} else if len(users) > 1 {
		return nil, fmt.Errorf("CHUserRep.GetByID: %w", ErrExpectedOneUser)
	}
	return users[0], nil
}

func (ch *CHUserRep) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, email, subscribeMail FROM Users WHERE login = ?"
	rows, err := ch.db.QueryContext(ctx, query, login)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetByLogin: %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	users, err := ch.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.GetByLogin: %v", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	} else if len(users) > 1 {
		return nil, fmt.Errorf("CHUserRep.GetByLogin: %w", ErrExpectedOneUser)
	}
	return users[0], nil
}

func (ch *CHUserRep) Add(ctx context.Context, u *models.User) error {
	_, err := ch.GetByLogin(ctx, u.GetLogin())
	if err == nil {
		return ErrDuplicateLoginUser
	} else if err != ErrUserNotFound {
		return fmt.Errorf("CHUserRep.Add %v", err)
	}

	subscribeMail := uint8(0)
	if u.IsSubscribedToMail() {
		subscribeMail = 1
	}

	var email interface{} = nil
	if u.GetEmail() != "" {
		email = u.GetEmail()
	}

	query := `INSERT INTO Users (id, username, login, hashedPassword, createdAt, email, subscribeMail) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := ch.db.ExecContext(ctx, query,
		u.GetID(),
		u.GetUsername(),
		u.GetLogin(),
		u.GetHashedPassword(),
		u.GetCreatedAt(),
		email,
		subscribeMail,
	)

	if err != nil {
		return fmt.Errorf("CHUserRep.Add %w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CHUserRep.Add %w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHUserRep) Delete(ctx context.Context, id uuid.UUID) error {
	query := "ALTER TABLE Users DELETE WHERE id = ?"
	result, err := ch.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("CHUserRep.Delete %w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CHUserRep.Delete %w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHUserRep) Update(
	ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.User) (*models.User, error),
) (*models.User, error) {
	user, err := ch.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.Update: %v", err)
	}

	updatedUser, err := funcUpdate(user)
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.Update funcUpdate: %v", err)
	}

	subscribeMail := uint8(0)
	if updatedUser.IsSubscribedToMail() {
		subscribeMail = 1
	}

	var email interface{} = nil
	if updatedUser.GetEmail() != "" {
		email = updatedUser.GetEmail()
	}

	query := `ALTER TABLE Users UPDATE 
		username = ?, 
		login = ?, 
		hashedPassword = ?, 
		email = ?, 
		subscribeMail = ? 
		WHERE id = ?`

	result, err := ch.db.ExecContext(ctx, query,
		updatedUser.GetUsername(),
		updatedUser.GetLogin(),
		updatedUser.GetHashedPassword(),
		email,
		subscribeMail,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("CHUserRep.Update: %w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("CHUserRep.Update: %w: %v", ErrRowsAffected, err)
	}
	return updatedUser, nil
}

func (ch *CHUserRep) UpdateSubscribeToMailing(ctx context.Context, id uuid.UUID, newSubscribeMail bool) error {
	subscribeMail := uint8(0)
	if newSubscribeMail {
		subscribeMail = 1
	}

	query := "ALTER TABLE Users UPDATE subscribeMail = ? WHERE id = ?"
	result, err := ch.db.ExecContext(ctx, query, subscribeMail, id)
	if err != nil {
		return fmt.Errorf("CHUserRep.UpdateSubscribeToMailing: %w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CHUserRep.UpdateSubscribeToMailing: %w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHUserRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHUserRep) Close() {
	ch.db.Close()
}
