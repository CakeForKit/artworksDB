package userrep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgUserRep struct {
	db *sql.DB
}

var (
	ErrOpenConnect = errors.New("open connect failed")
	ErrPing        = errors.New("ping failed")
)

func NewPgUserRep(ctx context.Context, pgCreds *cnfg.DatebaseCredentials, dbConf *cnfg.DatebaseConfig) (*PgUserRep, error) {
	// connStr := "postgres://puser:ppassword@postgres_artworks:5432/artworks"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("NewPgUserRep: %w: %w", ErrOpenConnect, err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("NewPgUserRep: %w: %w", ErrPing, err)
	}
	// Настраиваем пул соединений
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

	return &PgUserRep{db: db}, nil
}

func (pg *PgUserRep) parseUsersRows(rows *sql.Rows) ([]*models.User, error) {
	var resUsers []*models.User
	for rows.Next() {
		var id uuid.UUID
		var username, login, hashedPassword, email string
		var createdAt time.Time
		var subscribeMail bool
		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &email, &subscribeMail); err != nil {
			return nil, fmt.Errorf("parseUsersRows scan error: %w", err)
		}
		user, err := models.NewUser(id, username, login, hashedPassword, createdAt, email, subscribeMail)
		if err != nil {
			return nil, fmt.Errorf("parseUsersRows: %w", err)
		}
		resUsers = append(resUsers, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("parseUsersRows rows iteration error: %w", err)
	}
	return resUsers, nil
}

func (pg *PgUserRep) GetAll(ctx context.Context) ([]*models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetAll: %w: %w", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetAll: %w: %w", ErrQueryExec, err)
	}
	defer rows.Close()

	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetAll: %w", err)
	}
	return users, nil
}

func (pg *PgUserRep) GetAllSubscribed(ctx context.Context) ([]*models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
		From("users").
		Where(sq.Eq{"subscribeMail": true}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetAllSubscribed: %w: %w", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetAllSubscribed: %w: %w", ErrQueryExec, err)
	}
	defer rows.Close()

	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetAllSubscribed: %w", err)
	}

	return users, nil
}

func (pg *PgUserRep) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetByID: %w: %w", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetByID: %w: %w", ErrQueryExec, err)
	}
	defer rows.Close()
	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetByID: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	} else if len(users) > 1 {
		return nil, fmt.Errorf("PgUserRep.GetByID: %w: %w", ErrExpectedOneUser, err)
	}
	return users[0], nil
}

func (pg *PgUserRep) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
		From("users").
		Where(sq.Eq{"login": login}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetByLogin: %w: %w", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetByLogin: %w: %w", ErrQueryExec, err)
	}
	defer rows.Close()
	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.GetByLogin: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	} else if len(users) > 1 {
		return nil, fmt.Errorf("PgUserRep.GetByLogin: %w: %w", ErrExpectedOneUser, err)
	}
	return users[0], nil
}

func (pg *PgUserRep) Add(ctx context.Context, e *models.User) error {
	_, err := pg.GetByLogin(ctx, e.GetLogin())
	if err == nil {
		return ErrDuplicateLoginUser
	} else if err != ErrUserNotFound {
		return fmt.Errorf("PgUserRep.Add %w", err)
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert("Users").
		Columns("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
		Values(e.GetID(), e.GetUsername(), e.GetLogin(), e.GetHashedPassword(), e.GetCreatedAt(), e.GetEmail(), e.IsSubscribedToMail()).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgUserRep.Add %w: %w", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgUserRep.Add %w: %w", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgUserRep.Add %w: %w", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgUserRep.Add %w: no user added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgUserRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgUserRep.Delete %w: %w", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgUserRep.Delete %w: %w", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgUserRep.Delete %w: %w", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgUserRep.Delete %w: no user with id %s", ErrRowsAffected, id)
	}
	return nil
}

func (pg *PgUserRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.User) (*models.User, error)) (*models.User, error) {
	user, err := pg.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.Update: %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedUser, err := funcUpdate(user)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.Update funcUpdate: %w", err)
	}
	query, args, err := psql.Update("Users").
		Set("username", updatedUser.GetUsername()).
		Set("login", updatedUser.GetLogin()).
		Set("hashedPassword", updatedUser.GetHashedPassword()).
		Set("email", updatedUser.GetEmail()).
		Set("subscribeMail", updatedUser.IsSubscribedToMail()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.Update: %w: %w", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.Update: %w: %w", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("PgUserRep.Update: %w: %w", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("PgUserRep.Update: %w: no user updated", ErrRowsAffected)
	}
	return updatedUser, nil
}

func (pg *PgUserRep) UpdateSubscribeToMailing(ctx context.Context, id uuid.UUID, newSubscribeMail bool) error {
	_, err := pg.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("PgUserRep.Update: %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Update("Users").
		Set("subscribeMail", newSubscribeMail).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("PgUserRep.Update: %w: %w", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgUserRep.Update: %w: %w", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgUserRep.Update: %w: %w", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgUserRep.Update: %w: no user updated", ErrRowsAffected)
	}
	return nil
}

func (pg *PgUserRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgUserRep) Close() {
	if pg.db != nil {
		pg.db.Close()
	}
}
