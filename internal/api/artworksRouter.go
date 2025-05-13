package api

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArtworksRouter struct {
	artworksServ artworkserv.ArtworkService
}

func NewArtworksRouter(router *gin.RouterGroup, artworksServ artworkserv.ArtworkService) ArtworksRouter {
	r := ArtworksRouter{
		artworksServ: artworksServ,
	}
	gr := router.Group("artworks")
	gr.GET("/", r.GetAllArtworks)
	gr.POST("/", r.AddArtwork)
	gr.PUT("/", r.UpdateArtwork)
	gr.DELETE("/", r.DeleteArtwork)
	return r
}

// getAllArtworks godoc
// @Summary Get all artworks by employee
// @Description Retrieves a list of all artworks
// @Tags Artworks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.ArtworkResponse
// @Router /employee/artworks [get]
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
// @Description Add artwork with already created author and collection.
// @Tags Artworks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.AddArtworkRequest true "New Artwork with already created author and collection."
// @Success 201 "Artworks added"
// @Failure 400 "Wrong input parameters"
// @Failure 404 "Not Found"
// @Router /employee/artworks [post]
func (r *ArtworksRouter) AddArtwork(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.AddArtworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.artworksServ.Add(ctx, req)
	if err != nil {
		if errors.Is(err, authorrep.ErrAuthorNotFound) || errors.Is(err, collectionrep.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, models.ErrValidateArtwork) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// Update Artwork godoc
// @Summary Update artwork by employee
// @Description Update artwork with already created author and collection.
// @Tags Artworks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateArtworkRequest true "Updating Artwork with [new] author and [new] collection"
// @Success 200 "Artwork updated"
// @Failure 400 "Wrong input parameters"
// @Failure 404 "Not Found"
// @Router /employee/artworks [put]
func (r *ArtworksRouter) UpdateArtwork(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.UpdateArtworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.artworksServ.Update(
		ctx,
		uuid.MustParse(req.ID),
		jsonreqresp.ArtworkUpdate{
			Title:        req.Title,
			CreationYear: req.CreationYear,
			Technic:      req.Technic,
			Material:     req.Material,
			Size:         req.Size,
			AuthorID:     req.AuthorID,
			CollectionID: req.CollectionID,
		})

	if err != nil {
		if errors.Is(err, artworkrep.ErrArtworkNotFound) ||
			errors.Is(err, authorrep.ErrAuthorNotFound) ||
			errors.Is(err, collectionrep.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, artworkrep.ErrUpdateArtwork) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteArtwork godoc
// @Summary Delete an artwork
// @Description Deletes an existing artwork
// @Tags Artworks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.DeleteArtworkRequest true "Artwork delete data"
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Router /employee/artworks [delete]
func (r *ArtworksRouter) DeleteArtwork(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.DeleteArtworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.artworksServ.Delete(ctx, uuid.MustParse(req.ID))
	if err != nil {
		if errors.Is(err, artworkrep.ErrArtworkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
