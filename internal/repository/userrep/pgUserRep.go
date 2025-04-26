package userrep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/config"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgUserRep struct {
	db *sql.DB
}

var (
	pgInstance *PgUserRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect = errors.New("open connect failed")
	ErrPing        = errors.New("ping failed")
	ErrQueryBuilds = errors.New("query build failed")
	ErrQueryExec   = errors.New("query execution failed")
)

// func NewPgUserRep(ctx context.Context) (UserRep, error) {
func NewPgUserRep(ctx context.Context, conf *config.Config) (*PgUserRep, error) {
	var resErr error
	pgOnce.Do(func() {
		// connStr := "postgres://puser:ppassword@postgres_artworks:5432/artworks"
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			conf.Postgres.Username, conf.Postgres.Password, conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.DbName)
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
		db.SetMaxOpenConns(conf.Datebase.MaxIdleConns)
		db.SetMaxIdleConns(conf.Datebase.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(conf.Datebase.ConnMaxLifetime.Hours()))

		pgInstance = &PgUserRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

// func (pg *PgUserRep) TestSelect(ctx context.Context) error {
// 	if pg == nil || pg.db == nil {
// 		return errors.New("database connection is not initialized")
// 	}
// 	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
// 	query, args, err := psql.Select("username", "email").
// 		From("users").
// 		ToSql()
// 	if err != nil {
// 		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
// 	}

// 	rows, err := pg.db.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return fmt.Errorf("%w: %v", ErrQueryExec, err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var username, email string
// 		if err := rows.Scan(&username, &email); err != nil {
// 			return fmt.Errorf("scan error: %v", err)
// 		}
// 		// log.Printf("User: %s, Email: %s", username, email)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return fmt.Errorf("rows iteration error: %v", err)
// 	}
// 	return nil
// }

func (pg *PgUserRep) GetAll(ctx context.Context) ([]*models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	var resUsers []*models.User
	for rows.Next() {
		var id uuid.UUID
		var username, login, hashedPassword, email string
		var createdAt time.Time
		var subscribeMail bool
		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &email, &subscribeMail); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		user, err := models.NewUser(id, username, login, hashedPassword, createdAt, email, subscribeMail)
		if err != nil {
			return nil, err
		}
		resUsers = append(resUsers, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}
	return resUsers, nil
}

// func (pg *PgUserRep) GetAllSubscribed(ctx context.Context) ([]*models.User, error) {

// }

func (pg *PgUserRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgUserRep) Close() {
	pg.db.Close()
}
