package jsonreqresp

import "time"

type UserResponse struct {
	ID            string    `json:"id" example:"bb2e8400-e29b-41d4-a716-446655440000"`
	Username      string    `json:"username" example:"alice_smith"`
	Login         string    `json:"login" example:"alice@example.com"`
	CreatedAt     time.Time `json:"createdAt" example:"2023-06-15T14:30:00Z"`
	Email         string    `json:"email" example:"alice.smith@example.com"`
	SubscribeMail bool      `json:"subscribeMail" example:"true"`
}

type UserSelfResponse struct {
	Username      string `json:"username" example:"alice_smith"`
	Login         string `json:"login" example:"alice@example.com"`
	Email         string `json:"email" example:"alice.smith@example.com"`
	SubscribeMail bool   `json:"subscribeMail" example:"true"`
}
