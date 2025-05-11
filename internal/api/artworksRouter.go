package api

import (
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArtworksRouter struct {
	artworksServ artworkserv.ArtworkService
}

func (r *ArtworksRouter) Init(router *gin.RouterGroup, artworksServ artworkserv.ArtworkService) {
	r.artworksServ = artworksServ
	gr := router.Group("artworks")
	gr.GET("/all", r.GetAllArtworks)
	gr.POST("/add", r.AddArtwork)
	gr.PUT("/update", r.UpdateArtwork)
}

// getAllArtworks godoc
// @Summary Get all artworks by employee
// @Description Retrieves a list of all artworks
// @Tags employee
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.ArtworkResponse
// @Router /employee/artworks/all [get]
func (r *ArtworksRouter) GetAllArtworks(c *gin.Context) {
	ctx := c.Request.Context()
	artworks, err := r.artworksServ.GetAllArtworks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	artworksResp := make([]jsonreqresp.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}
	c.JSON(http.StatusOK, artworksResp)
}

// AddArtwork godoc
// @Summary Add artwork by employee
// @Description Add artwork with [new] author and [new] collection. If author or collection ID = "", an it will be created.
// @Tags employee
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.ArtworkRequest true "New Artwork with [new] author and [new] collection"
// @Success 201 "Artworks added"
// @Failure 400 "Wrong input parameters"
// @Router /employee/artworks/add [post]
func (r *ArtworksRouter) AddArtwork(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.ArtworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	artwork, err := models.FromArtworkRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.artworksServ.Add(ctx, &artwork)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// Update Artwork godoc
// @Summary Update artwork by employee
// @Description Update artwork with [new] author and [new] collection. If author or collection ID = "", an it will be created.
// @Tags employee
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateArtworkRequest true "Updating Artwork with [new] author and [new] collection"
// @Success 200 "Artwork updated"
// @Failure 400 "Wrong input parameters"
// @Router /employee/artworks/update [put]
func (r *ArtworksRouter) UpdateArtwork(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.UpdateArtworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := models.FromAuthorRequest(req.Author)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection, err := models.FromCollectionRequest(req.Collection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	artworkID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	funcUpdate := func(a *models.Artwork) (*models.Artwork, error) {
		art, err := models.NewArtwork(
			artworkID,
			req.Title,
			req.Technic,
			req.Material,
			req.Size,
			req.CreationYear,
			&author,
			&collection,
		)
		return &art, err
	}

	err = r.artworksServ.Update(ctx, artworkID, funcUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// func (r *ArtworksRouter) DeleteArtwork(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	var req jsonreqresp.DeleteArtworkRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	artworkID, err := uuid.Parse(req.ID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err = r.artworksServ.Delete(ctx, artworkID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// }
