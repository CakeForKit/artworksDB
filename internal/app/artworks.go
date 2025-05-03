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
