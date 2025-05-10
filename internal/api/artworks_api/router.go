package artworksapi

import (
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"github.com/gin-gonic/gin"
)

type ArtworksRouter struct {
	artworksServ artworkserv.ArtworkService
}

func (r *ArtworksRouter) Init(router *gin.RouterGroup, artworksServ artworkserv.ArtworkService) {
	r.artworksServ = artworksServ
	gr := router.Group("artworks")
	gr.GET("/all", r.getAllArtworks)
}

// getAllArtworks godoc
// @Summary Get all artworks by employee
// @Description Retrieves a list of all artworks
// @Tags employee
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} models.ArtworkResponse
// @Failure 401 "Unauthorized"
// @Router /employee/artworks/all [get]
func (r *ArtworksRouter) getAllArtworks(c *gin.Context) {
	ctx := c.Request.Context()
	artworks, err := r.artworksServ.GetAllArtworks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	artworksResp := make([]models.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}
	c.JSON(http.StatusOK, artworksResp)
}
