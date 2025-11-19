package jsonreqresp

type AuthorResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"name" example:"Leonardo da Vinci"`
	BirthYear int    `json:"birthYear" example:"1452"`
	DeathYear int    `json:"deathYear" example:"1519"`
}

type AddAuthorRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=100" example:"Винсент Ван Гог"`
	BirthYear int    `json:"birthYear" binding:"required,gte=1000" example:"1853"`
	DeathYear int    `json:"deathYear,omitempty" binding:"omitempty,gtefield=BirthYear" example:"1890"`
}

type UpdateAuthorRequest struct {
	ID        string `json:"id" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
	Name      string `json:"name" binding:"required,min=2,max=100" example:"Винсент Ван Гог"`
	BirthYear int    `json:"birthYear" binding:"required,gte=1000" example:"1853"`
	DeathYear int    `json:"deathYear,omitempty" binding:"omitempty,gtefield=BirthYear" example:"1890"`
}

type DeleteAuthorRequest struct {
	ID string `json:"id" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
}
