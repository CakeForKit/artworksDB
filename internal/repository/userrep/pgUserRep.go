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
	ErrOpenConnect     = errors.New("open connect failed")
	ErrPing            = errors.New("ping failed")
	ErrQueryBuilds     = errors.New("query build failed")
	ErrQueryExec       = errors.New("query execution failed")
	ErrNoUser          = errors.New("no user")
	ErrExpectedOneUser = errors.New("expected one user")
	ErrRowsAffected    = errors.New("no rows affected")
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

func (pg *PgUserRep) parseUsersRows(rows *sql.Rows) ([]*models.User, error) {
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

	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("%w: %v", ErrNoUser, err)
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
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("%w: %v", ErrNoUser, err)
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
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("%w: %v", ErrNoUser, err)
	} else if len(users) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneUser, err)
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
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("%w: %v", ErrNoUser, err)
	} else if len(users) > 1 {
		return nil, fmt.Errorf("%w: %v", ErrExpectedOneUser, err)
	}
	return users[0], nil
}

func (pg *PgUserRep) Add(ctx context.Context, e *models.User) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert("Users").
		Columns("id", "username", "login", "hashedPassword", "createdAt", "email", "subscribeMail").
		Values(e.GetID(), e.GetUsername(), e.GetLogin(), e.GetHashedPassword(), e.GetCreatedAt(), e.GetMail(), e.IsSubscribedToMail()).
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
		return fmt.Errorf("%w: no user added", ErrNoUser)
	}
	return nil
}

func (pg *PgUserRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Users").
		Where(sq.Eq{"id": id}).
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
		return fmt.Errorf("%w: no user with id %s", ErrNoUser, id)
	}
	return nil
}

func (pg *PgUserRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.User) (*models.User, error)) (*models.User, error) {
	user, err := pg.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedUser, err := funcUpdate(user)
	if err != nil {
		return nil, fmt.Errorf("funcUpdate: %v", err)
	}
	query, args, err := psql.Update("Users").
		Set("username", updatedUser.GetUsername()).
		Set("login", updatedUser.GetLogin()).
		Set("hashedPassword", updatedUser.GetHashedPassword()).
		Set("email", updatedUser.GetMail()).
		Set("subscribeMail", updatedUser.IsSubscribedToMail()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("%w: no user added", ErrNoUser)
	}
	return updatedUser, nil
}

func (pg *PgUserRep) UpdateSubscribeToMailing(ctx context.Context, id uuid.UUID, newSubscribeMail bool) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Update("Users").
		Set("subscribeMail", newSubscribeMail).
		Where(sq.Eq{"id": id}).ToSql()
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
		return fmt.Errorf("%w: no user added", ErrNoUser)
	}
	return nil
}

func (pg *PgUserRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgUserRep) Close() {
	pg.db.Close()
}
