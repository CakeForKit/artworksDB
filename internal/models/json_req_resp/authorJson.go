package jsonreqresp

type AuthorResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"name" example:"Leonardo da Vinci"`
	BirthYear int    `json:"birthYear" example:"1452"`
	DeathYear int    `json:"deathYear" example:"1519"`
}

type AuthorRequest struct {
	ID        string `json:"id,omitempty" example:"ba1df957-ed5e-4694-8766-c5ec5806e5e7"`
	Name      string `json:"name" binding:"required,min=2,max=100" example:"Винсент Ван Гог"`           // Обязательное, 2-100 символов
	BirthYear int    `json:"birthYear" binding:"required,gte=1000" example:"1853"`                      // Обязательное, >= 1000
	DeathYear int    `json:"deathYear,omitempty" binding:"omitempty,gtefield=BirthYear" example:"1890"` // Опциональное, >= BirthYear
}
