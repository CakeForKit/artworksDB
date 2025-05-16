package employeerep

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

type PgEmployeeRep struct {
	db *sql.DB
}

var (
	pgInstance *PgEmployeeRep
	pgOnce     sync.Once
)

var (
	ErrOpenConnect         = errors.New("open connect failed")
	ErrPing                = errors.New("ping failed")
	ErrQueryBuilds         = errors.New("query build failed")
	ErrQueryExec           = errors.New("query execution failed")
	ErrExpectedOneEmployee = errors.New("expected one employee")
	ErrRowsAffected        = errors.New("no rows affected")
)

func NewPgEmployeeRep(ctx context.Context, pgCreds *cnfg.PostgresCredentials, dbConf *cnfg.DatebaseConfig) (*PgEmployeeRep, error) {
	var resErr error
	pgOnce.Do(func() {
		// connStr := "postgres://puser:ppassword@postgres_artworks:5432/artworks"
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			resErr = fmt.Errorf("NewPgEmployeeRep: %w: %v", ErrOpenConnect, err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			resErr = fmt.Errorf("NewPgEmployeeRep: %w: %v", ErrPing, err)
			db.Close()
			return
		}
		// Настраиваем пул соединений
		db.SetMaxOpenConns(dbConf.MaxOpenConns)
		db.SetMaxIdleConns(dbConf.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		pgInstance = &PgEmployeeRep{db: db}
	})
	if resErr != nil {
		return nil, resErr
	}

	return pgInstance, nil
}

func (pg *PgEmployeeRep) parseEmployeesRows(rows *sql.Rows) ([]*models.Employee, error) {
	var resEmployees []*models.Employee
	for rows.Next() {
		var id, adminID uuid.UUID
		var username, login, hashedPassword string
		var createdAt time.Time
		var valid bool
		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &valid, &adminID); err != nil {
			return nil, fmt.Errorf("parseEmployeesRows: scan error: %v", err)
		}
		employee, err := models.NewEmployee(id, username, login, hashedPassword, createdAt, valid, adminID)
		if err != nil {
			return nil, err
		}
		resEmployees = append(resEmployees, &employee)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("parseEmployeesRows: rows iteration error: %v", err)
	}
	return resEmployees, nil
}

func (pg *PgEmployeeRep) GetAll(ctx context.Context) ([]*models.Employee, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "valid", "adminID").
		From("employees").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetAll %w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetAll %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	employees, err := pg.parseEmployeesRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetAll %v", err)
	}
	if len(employees) == 0 {
		return nil, ErrEmployeeNotFound
	}
	return employees, nil
}

func (pg *PgEmployeeRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "valid", "adminID").
		From("Employees").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetByID %w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetByID %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	employees, err := pg.parseEmployeesRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetByID %v", err)
	}
	if len(employees) == 0 {
		return nil, ErrEmployeeNotFound
	} else if len(employees) > 1 {
		return nil, fmt.Errorf("PgEmployeeRep.GetByID %w: %v", ErrExpectedOneEmployee, err)
	}
	return employees[0], nil
}

func (pg *PgEmployeeRep) GetByLogin(ctx context.Context, login string) (*models.Employee, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "username", "login", "hashedPassword", "createdAt", "valid", "adminID").
		From("employees").
		Where(sq.Eq{"login": login}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetByLogin %w: %v", ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetByLogin %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()
	employees, err := pg.parseEmployeesRows(rows)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.GetByLogin %v", err)
	}
	if len(employees) == 0 {
		return nil, ErrEmployeeNotFound
	} else if len(employees) > 1 {
		return nil, fmt.Errorf("PgEmployeeRep.GetByLogin %w: %v", ErrExpectedOneEmployee, err)
	}
	return employees[0], nil
}

func (pg *PgEmployeeRep) Add(ctx context.Context, e *models.Employee) error {
	_, err := pg.GetByLogin(ctx, e.GetLogin())
	if err == nil {
		return ErrDuplicateLoginEmp
	} else if err != ErrEmployeeNotFound {
		return fmt.Errorf("PgEmployeeRep.Add %v", err)
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert("Employees").
		Columns("id", "username", "login", "hashedPassword", "createdAt", "valid", "adminID").
		Values(e.GetID(), e.GetUsername(), e.GetLogin(), e.GetHashedPassword(), e.GetCreatedAt(), true, e.GetAdminID()).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgEmployeeRep.Add %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgEmployeeRep.Add %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgEmployeeRep.Add %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgEmployeeRep.Add %w: no employee added", ErrRowsAffected)
	}
	return nil
}

func (pg *PgEmployeeRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Employees").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgEmployeeRep.Delete %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PgEmployeeRep.Delete %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PgEmployeeRep.Delete %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("PgEmployeeRep.Delete %w: no employee with id %s", ErrRowsAffected, id)
	}
	return nil
}

func (pg *PgEmployeeRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Employee) (*models.Employee, error)) (*models.Employee, error) {
	employee, err := pg.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.Update %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedEmployee, err := funcUpdate(employee)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.Update funcUpdate: %v", err)
	}
	query, args, err := psql.Update("Employees").
		Set("username", updatedEmployee.GetUsername()).
		Set("login", updatedEmployee.GetLogin()).
		Set("hashedPassword", updatedEmployee.GetHashedPassword()).
		Set("valid", updatedEmployee.IsValid()).
		Set("adminID", updatedEmployee.GetAdminID()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.Update %w: %v", ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.Update %w: %v", ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("PgEmployeeRep.Update %w: %v", ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("PgEmployeeRep.Update %w: no employee added", ErrRowsAffected)
	}
	return updatedEmployee, nil
}

func (pg *PgEmployeeRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgEmployeeRep) Close() {
	pg.db.Close()
}
