package jsonreqresp

type ArtworkResponse struct {
	ID           string             `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title        string             `json:"title" example:"Mona Lisa"`
	CreationYear int                `json:"creationYear" example:"1503"`
	Technic      string             `json:"technic" example:"Oil painting"`
	Material     string             `json:"material" example:"Poplar wood"`
	Size         string             `json:"size" example:"77 cm × 53 cm"`
	Author       AuthorResponse     `json:"author"`
	Collection   CollectionResponse `json:"collection"`
}

type AddArtworkRequest struct {
	Title        string `json:"title" binding:"required,max=255" example:"Звёздная ночь"`
	CreationYear int    `json:"creationYear" binding:"required,gt=0,lte=2100" example:"1889"`
	Technic      string `json:"technic" binding:"required,max=100" example:"Масло, холст"`
	Material     string `json:"material" binding:"required,max=100" example:"Холст, масляные краски"`
	Size         string `json:"size" binding:"required,max=50" example:"73.7 × 92.1 см"`
	AuthorID     string `json:"authorID" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
	CollectionID string `json:"collectionID" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
}

type ArtworkUpdate struct {
	Title        string `json:"title" binding:"required,max=255" example:"Звёздная ночь"`
	CreationYear int    `json:"creationYear" binding:"required,gt=0,lte=2100" example:"1889"`
	Technic      string `json:"technic" binding:"required,max=100" example:"Масло, холст"`
	Material     string `json:"material" binding:"required,max=100" example:"Холст, масляные краски"`
	Size         string `json:"size" binding:"required,max=50" example:"73.7 × 92.1 см"`
	AuthorID     string `json:"authorID" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
	CollectionID string `json:"collectionID" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
}

type UpdateArtworkRequest struct {
	ID           string `json:"id" binding:"required,uuid" example: "44a315d0-663c-4813-92a6-d7977c2f2aba"`
	Title        string `json:"title" binding:"required,max=255" example:"Звёздная ночь"`
	CreationYear int    `json:"creationYear" binding:"required,gt=0,lte=2100" example:"1889"`
	Technic      string `json:"technic" binding:"required,max=100" example:"Масло, холст"`
	Material     string `json:"material" binding:"required,max=100" example:"Холст, масляные краски"`
	Size         string `json:"size" binding:"required,max=50" example:"73.7 × 92.1 см"`
	AuthorID     string `json:"authorID" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
	CollectionID string `json:"collectionID" binding:"required,uuid" example:"cfd9ff5d-cb37-407c-b043-288a482e9239"`
}

type DeleteArtworkRequest struct {
	ID string `json:"id" binding:"required,uuid"`
}
