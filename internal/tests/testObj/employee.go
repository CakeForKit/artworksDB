package testobj

import (
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/google/uuid"
)

type EmployeeMother interface {
	DefaultEmployee(employeeID uuid.UUID, adminID uuid.UUID) models.Employee
	EmployeeWithPswdHash(employeeID uuid.UUID, adminID uuid.UUID, hashedPassword string) models.Employee
	DefaultEmployeeP(employeeID uuid.UUID, adminID uuid.UUID) *models.Employee
	EmployeeWithLoginP(employeeID uuid.UUID, adminID uuid.UUID, login string) *models.Employee
}

func NewEmployeeMother() EmployeeMother {
	return &employeeMother{}
}

type employeeMother struct{}

func (um *employeeMother) DefaultEmployee(employeeID uuid.UUID, adminID uuid.UUID) models.Employee {
	employee, _ := models.NewEmployee(
		employeeID,
		"test-employee",
		"test-login"+employeeID.String(),
		"hashed-password",
		time.Now(),
		true,
		adminID,
	)
	return employee
}

func (um *employeeMother) EmployeeWithPswdHash(
	employeeID uuid.UUID, adminID uuid.UUID, hashedPassword string) models.Employee {
	employee, _ := models.NewEmployee(
		employeeID,
		"test-employee",
		"test-login"+employeeID.String(),
		hashedPassword,
		time.Now(),
		true,
		adminID,
	)
	return employee
}

func (um *employeeMother) DefaultEmployeeP(employeeID uuid.UUID, adminID uuid.UUID) *models.Employee {
	employee, _ := models.NewEmployee(
		employeeID,
		"test-employee",
		"test-login"+employeeID.String(),
		"hashedpassword",
		time.Now(),
		true,
		adminID,
	)
	return &employee
}

func (um *employeeMother) EmployeeWithLoginP(employeeID uuid.UUID, adminID uuid.UUID, login string) *models.Employee {
	employee, _ := models.NewEmployee(
		employeeID,
		"test-employee",
		login,
		"hashedpassword",
		time.Now(),
		true,
		adminID,
	)
	return &employee
}
