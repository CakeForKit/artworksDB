package adminrep

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

type PgAdminRep struct {
	db *sql.DB
}

var (
	pgInstance *PgAdminRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect      = errors.New("open connect failed")
	ErrPing             = errors.New("ping failed")
	ErrQueryBuilds      = errors.New("query build failed")
	ErrQueryExec        = errors.New("query execution failed")
	ErrExpectedOneAdmin = errors.New("expected one admin")
	ErrRowsAffected     = errors.New("no rows affected")
)

func NewPgAdminRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (*PgAdminRep, error) {
	var resErr error
	pgOnce.Do(func() {
		// connStr := "postgres://puser:ppassword@postgres_artworks:5432/artworks"
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			resErr = fmt.Errorf("NewPgAdminRep: %w: %v", ErrOpenConnect, err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			resErr = fmt.Errorf("NewPgAdminRep: %w: %v", ErrPing, err)
			db.Close()
			return
		}
		// Настраиваем пул соединений
		db.SetMaxOpenConns(dbConf.MaxOpenConns)
		db.SetMaxIdleConns(dbConf.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		pgInstance = &PgAdminRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

func (pg *PgAdminRep) parseAdminsRows(rows *sql.Rows) ([]*models.Admin, error) {
	var resAdmins []*models.Admin
	for rows.Next() {
		var id uuid.UUID
		var username, login, hashedPassword string
		var createdAt time.Time
		var valid bool
		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &valid); err != nil {
			return nil, fmt.Errorf("parseAdminsRows, scan error: %v", err)
		}
		admin, err := models.NewAdmin(id, username, login, hashedPassword, createdAt, valid)
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

func (pg *PgAdminRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Admin, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	res, err := pg.parseAdminsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

func (pg *PgAdminRep) GetAll(ctx context.Context) ([]*models.Admin, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "valid").
		From("Admins")
	res, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgAdminRep.GetAll: %v", err)
	}
	return res, nil
}

func (pg *PgAdminRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Admin, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "valid").
		From("Admins").
		Where(sq.Eq{"id": id})
	res, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgAdminRep.GetByID: %v", err)
	}
	if len(res) == 0 {
		return nil, ErrAdminNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("PgAdminRep.GetByID: %w: %v", ErrExpectedOneAdmin, err)
	}
	return res[0], nil
}

func (pg *PgAdminRep) GetByLogin(ctx context.Context, login string) (*models.Admin, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "valid").
		From("Admins").
		Where(sq.Eq{"login": login})
	res, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgAdminRep.GetByID: %v", err)
	}
	if len(res) == 0 {
		return nil, ErrAdminNotFound
	} else if len(res) > 1 {
		return nil, fmt.Errorf("PgAdminRep.GetByLogin: %w: %v", ErrExpectedOneAdmin, err)
	}
	return res[0], nil
}

func (pg *PgAdminRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, querySQL, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: no added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgAdminRep) Add(ctx context.Context, e *models.Admin) error {
	_, err := pg.GetByLogin(ctx, e.GetLogin())
	if err == nil {
		return ErrDuplicateLoginAdm
	} else if err != ErrAdminNotFound {
		return fmt.Errorf("PgAdminRep.Add %v", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("admins").
		Columns("id", "username", "login", "hashedPassword", "createdAt", "valid").
		Values(e.GetID(), e.GetUsername(), e.GetLogin(), e.GetHashedPassword(), e.GetCreatedAt(), true)
	err = pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgAdminRep.Add: %w", err)
	}
	return nil
}

func (pg *PgAdminRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("Admins").
		Where(sq.Eq{"id": id})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgAdminRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgAdminRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Admin) (*models.Admin, error),
) error {
	admin, err := pg.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("PgAdminRep.Update %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedAdmin, err := funcUpdate(admin)
	if err != nil {
		return fmt.Errorf("PgAdminRep.Update: %w", ErrUpdateAdmin)
	}
	query := psql.Update("Admins").
		Set("username", updatedAdmin.GetUsername()).
		Set("login", updatedAdmin.GetLogin()).
		Set("hashedPassword", updatedAdmin.GetHashedPassword()).
		Set("valid", updatedAdmin.IsValid()).
		Where(sq.Eq{"id": id})
	err = pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgAdminRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgAdminRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgAdminRep) Close() {
	pg.db.Close()
}
