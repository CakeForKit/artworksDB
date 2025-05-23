package jsonreqresp

import "time"

type EmployeeResponse struct {
	ID        string    `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Username  string    `json:"username" example:"john doe"`
	Login     string    `json:"login" example:"johndoe@example.com"`
	CreatedAt time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	Valid     bool      `json:"valid" example:"true"`
	AdminID   string    `json:"adminId" example:"bb2e8400-e29b-41d4-a716-446655443333"`
}

type UpdateValidEmployeeRequest struct {
	ID    string `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Valid bool   `json:"valid" example:"true"`
}

// type CreateEmployeeRequest struct {
// 	Username string `json:"username" binding:"required"`
// 	Login    string `json:"login" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=8"`
// 	AdminID  string `json:"adminId" binding:"required,uuid"`
// }

// type UpdateEmployeeRequest struct {
// 	Username *string `json:"username,omitempty"`
// 	Login    *string `json:"login,omitempty" validate:"omitempty,email"`
// 	Password *string `json:"password,omitempty" validate:"omitempty,min=8"`
// 	Valid    *bool   `json:"valid,omitempty"`
// 	AdminID  *string `json:"adminId,omitempty" validate:"omitempty,uuid"`
// }
