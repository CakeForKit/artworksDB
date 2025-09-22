package testobj

import (
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/google/uuid"
)

type EmployeeMother interface {
	DefaultEmployee(employeeID uuid.UUID, adminID uuid.UUID) models.Employee
	EmployeeWithPswdHash(employeeID uuid.UUID, adminID uuid.UUID, hashedPassword string) models.Employee
	DefaultEmployeeP(employeeID uuid.UUID, adminID uuid.UUID) *models.Employee
}

func NewEmployeeMother() EmployeeMother {
	return &employeeMother{}
}

type employeeMother struct{}

func (um *employeeMother) DefaultEmployee(employeeID uuid.UUID, adminID uuid.UUID) models.Employee {
	employee, _ := models.NewEmployee(
		employeeID,
		"test-employee",
		"test-login"+uuid.NewString(),
		"hashed-password",
		time.Now(),
		true,
		adminID,
	)
	return employee
}

func (um *employeeMother) EmployeeWithPswdHash(employeeID uuid.UUID, adminID uuid.UUID, hashedPassword string) models.Employee {
	employee, _ := models.NewEmployee(
		employeeID,
		"test-employee",
		"test-login"+uuid.NewString(),
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
		"test-login"+uuid.NewString(),
		"hashedpassword",
		time.Now(),
		true,
		adminID,
	)
	return &employee
}
