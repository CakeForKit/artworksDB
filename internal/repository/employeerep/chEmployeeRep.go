package employeerep

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

type CHEmployeeRep struct {
	db *sql.DB
}

var (
	chInstance *CHEmployeeRep
	chOnce     sync.Once
)

func NewCHEmployeeRep(ctx context.Context, chCreds *cnfg.ClickHouseCredentials, dbConf *cnfg.DatebaseConfig) (*CHEmployeeRep, error) {
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
			resErr = fmt.Errorf("NewCHEmployeeRep: %w: %v", ErrPing, err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(dbConf.MaxOpenConns)
		conn.SetMaxIdleConns(dbConf.MaxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

		chInstance = &CHEmployeeRep{db: conn}
	})
	if resErr != nil {
		return nil, resErr
	}

	return chInstance, nil
}

func (ch *CHEmployeeRep) parseEmployeesRows(rows *sql.Rows) ([]*models.Employee, error) {
	var resEmployees []*models.Employee
	for rows.Next() {
		var id, adminID uuid.UUID
		var username, login, hashedPassword string
		var createdAt time.Time
		var valid uint8

		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &valid, &adminID); err != nil {
			return nil, fmt.Errorf("parseEmployeesRows: scan error: %v", err)
		}

		employee, err := models.NewEmployee(id, username, login, hashedPassword, createdAt, valid == 1, adminID)
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

func (ch *CHEmployeeRep) GetAll(ctx context.Context) ([]*models.Employee, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, valid, adminID FROM Employees"
	rows, err := ch.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.GetAll %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	employees, err := ch.parseEmployeesRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.GetAll %v", err)
	}
	if len(employees) == 0 {
		return nil, ErrEmployeeNotFound
	}
	return employees, nil
}

func (ch *CHEmployeeRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, valid, adminID FROM Employees WHERE id = ?"
	rows, err := ch.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.GetByID %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	employees, err := ch.parseEmployeesRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.GetByID %v", err)
	}
	if len(employees) == 0 {
		return nil, ErrEmployeeNotFound
	} else if len(employees) > 1 {
		return nil, fmt.Errorf("CHEmployeeRep.GetByID %w: %v", ErrExpectedOneEmployee, err)
	}
	return employees[0], nil
}

func (ch *CHEmployeeRep) GetByLogin(ctx context.Context, login string) (*models.Employee, error) {
	query := "SELECT id, username, login, hashedPassword, createdAt, valid, adminID FROM Employees WHERE login = ?"
	rows, err := ch.db.QueryContext(ctx, query, login)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.GetByLogin %w: %v", ErrQueryExec, err)
	}
	defer rows.Close()

	employees, err := ch.parseEmployeesRows(rows)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.GetByLogin %v", err)
	}
	if len(employees) == 0 {
		return nil, ErrEmployeeNotFound
	} else if len(employees) > 1 {
		return nil, fmt.Errorf("CHEmployeeRep.GetByLogin %w: %v", ErrExpectedOneEmployee, err)
	}
	return employees[0], nil
}

func (ch *CHEmployeeRep) Add(ctx context.Context, e *models.Employee) error {
	_, err := ch.GetByLogin(ctx, e.GetLogin())
	if err == nil {
		return ErrDuplicateLoginEmp
	} else if err != ErrEmployeeNotFound {
		return fmt.Errorf("CHEmployeeRep.Add %v", err)
	}

	valid := uint8(0)
	if e.IsValid() {
		valid = 1
	}

	query := `
		INSERT INTO Employees 
		(id, username, login, hashedPassword, createdAt, valid, adminID) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := ch.db.ExecContext(ctx, query,
		e.GetID(),
		e.GetUsername(),
		e.GetLogin(),
		e.GetHashedPassword(),
		e.GetCreatedAt(),
		valid,
		e.GetAdminID(),
	)

	if err != nil {
		return fmt.Errorf("CHEmployeeRep.Add %w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CHEmployeeRep.Add %w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHEmployeeRep) Delete(ctx context.Context, id uuid.UUID) error {
	query := "ALTER TABLE Employees DELETE WHERE id = ?"
	result, err := ch.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("CHEmployeeRep.Delete %w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CHEmployeeRep.Delete %w: %v", ErrRowsAffected, err)
	}
	return nil
}

func (ch *CHEmployeeRep) Update(
	ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.Employee) (*models.Employee, error),
) (*models.Employee, error) {
	employee, err := ch.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.Update %w", err)
	}

	updatedEmployee, err := funcUpdate(employee)
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.Update funcUpdate: %v", err)
	}

	valid := uint8(0)
	if updatedEmployee.IsValid() {
		valid = 1
	}

	query := `
		ALTER TABLE Employees UPDATE 
		username = ?, 
		login = ?, 
		hashedPassword = ?, 
		valid = ?, 
		adminID = ? 
		WHERE id = ?`

	result, err := ch.db.ExecContext(ctx, query,
		updatedEmployee.GetUsername(),
		updatedEmployee.GetLogin(),
		updatedEmployee.GetHashedPassword(),
		valid,
		updatedEmployee.GetAdminID(),
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.Update %w: %v", ErrQueryExec, err)
	}

	// ClickHouse has limited RowsAffected support
	_, err = result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("CHEmployeeRep.Update %w: %v", ErrRowsAffected, err)
	}
	return updatedEmployee, nil
}

func (ch *CHEmployeeRep) Ping(ctx context.Context) error {
	return ch.db.PingContext(ctx)
}

func (ch *CHEmployeeRep) Close() {
	ch.db.Close()
}
