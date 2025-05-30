package jsonreqresp

type CollectionResponse struct {
	ID    string `json:"id" example:"aa1e8400-e29b-41d4-a716-446655441111"`
	Title string `json:"title" example:"Louvre Museum Collection"`
}

type AddCollectionRequest struct {
	Title string `json:"title" binding:"required,min=2,max=255" example:"Музей современного искусства"` // Обязательное, 2-255 символов
}

type UpdateCollectionRequest struct {
	ID    string `json:"id" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
	Title string `json:"title" binding:"required,min=2,max=255" example:"Музей современного искусства"` // Обязательное, 2-255 символов
}

type DeleteCollectionRequest struct {
	ID string `json:"id" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
}
