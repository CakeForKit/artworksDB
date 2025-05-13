package jsonreqresp

import "time"

type EventResponse struct {
	ID         string    `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title      string    `json:"title" example:"Выставка импрессионистов"`
	DateBegin  time.Time `json:"dateBegin" example:"2023-06-15T10:00:00Z"`
	DateEnd    time.Time `json:"dateEnd" example:"2023-09-20T18:00:00Z"`
	Address    string    `json:"address" example:"ул. Пречистенка, 12/2"`
	CanVisit   bool      `json:"canVisit" example:"true"`
	EmployeeID string    `json:"employeeID" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
	CntTickets int       `json:"cntTickets" example:"150"`
	Valid      bool      `json:"valid" example:"true"`
	ArtworkIDs []string  `json:"artworkIDs"`
}

type EventUpdate struct {
	Title      string    `json:"title" binding:"required,max=255" example:"Ночная выставка"`
	DateBegin  time.Time `json:"dateBegin" binding:"required" example:"2023-06-15T10:00:00Z"`
	DateEnd    time.Time `json:"dateEnd" binding:"required" example:"2023-09-20T18:00:00Z"`
	Address    string    `json:"address" binding:"required,max=500" example:"ул. Пречистенка, 12/2"`
	CanVisit   bool      `json:"canVisit" example:"true"`
	CntTickets int       `json:"cntTickets" binding:"required,min=0" example:"100"`
	Valid      bool      `json:"valid" example:"true"`
}

type AddEventRequest struct {
	Title      string    `json:"title" binding:"required,max=255" example:"Ночная выставка"`
	DateBegin  time.Time `json:"dateBegin" binding:"required" example:"2023-06-15T10:00:00Z"`
	DateEnd    time.Time `json:"dateEnd" binding:"required" example:"2023-09-20T18:00:00Z"`
	Address    string    `json:"address" binding:"required,max=500" example:"ул. Пречистенка, 12/2"`
	CanVisit   bool      `json:"canVisit" example:"true"`
	EmployeeID string    `json:"employeeID" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
	CntTickets int       `json:"cntTickets" binding:"required,min=0" example:"100"`
	ArtworkIDs []string  `json:"artworkIDs"`
}

type UpdateEventRequest struct {
	ID         string    `json:"id" binding:"required,uuid" example:"44a315d0-663c-4813-92a6-d7977c2f2aba"`
	Title      string    `json:"title" binding:"required,max=255" example:"Ночная выставка"`
	DateBegin  time.Time `json:"dateBegin" binding:"required" example:"2023-06-15T10:00:00Z"`
	DateEnd    time.Time `json:"dateEnd" binding:"required" example:"2023-09-20T18:00:00Z"`
	Address    string    `json:"address" binding:"required,max=500" example:"ул. Пречистенка, 12/2"`
	CanVisit   bool      `json:"canVisit" example:"true"`
	CntTickets int       `json:"cntTickets" binding:"required,min=0" example:"100"`
	Valid      bool      `json:"valid" example:"true"`
}

type DeleteEventRequest struct {
	ID string `json:"id" binding:"required,uuid"`
}

type ConArtworkEventRequest struct {
	EventID   string `json:"eventID" binding:"required,uuid"`
	ArtworkID string `json:"artworkID" binding:"required,uuid"`
}
