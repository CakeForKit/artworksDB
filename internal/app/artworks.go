package app

import (
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"github.com/gin-gonic/gin"
)

// getAllArtworks godoc
// @Summary Get all artworks
// @Description Retrieves a list of all artworks
// @Tags artworks
// @Accept json
// @Produce json
// @Success 200 {array} models.ArtworkResponse
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /artworks [get]
func (s *Server) getAllArtworks(ctx *gin.Context) {
	artworks, err := s.artworkServ.GetAllArtworks(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}
	artworksResp := make([]models.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}
	ctx.JSON(http.StatusOK, artworksResp)
}

// func (s *Server) AddAuthor(ctx *gin.Context) {
// 	var req models.AuthorRequest

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
// 		return
// 	}
// 	author, err := models.NewAuthor(
// 		uuid.New(),
// 		req.Name,
// 		req.BirthYear,
// 		*req.DeathYear,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
// 		return
// 	}

// 	s.artworkServ.Add(ctx.Request.Context(), &author)
// 	// Далее сохранение в БД...
// 	c.JSON(http.StatusCreated, author)
// }

// // @Router /artworks [get]
// func (s *Server) AddArtwork(ctx *gin.Context) error {
// 	var req CreateArtworkRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	artwork, err := ToArtworkModel(req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
// 		return
// 	}

// 	// Сохраняем artwork в БД...

// 	c.JSON(http.StatusCreated, ToArtworkResponse(artwork))
// }
