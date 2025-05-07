package querytime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type QueryTime interface {
	MeasureTime()
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

func NewQueryTime(
	ctx context.Context,
	pgCreds *cnfg.PostgresCredentials,
	dbConf *cnfg.DatebaseConfig,
) (QueryTime, error) {
	var resErr error
	pgOnceRep.Do(func() {
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

		pgInstance = &queryTime{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

// func (q *queryTime) createContaiter() {
// 	ctx := context.Background()
// 	dbCnfg := cnfg.GetTestDatebaseConfig()
// }

// func (q *queryTime) MeasureTime() {
// 	ctx := context.Background()
// 	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
// 	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
// 		From("users").
// 		ToSql()
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
// 	}
// }
