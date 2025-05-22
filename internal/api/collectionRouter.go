package api

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/collectionserv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CollectionRouter struct {
	collectionServ collectionserv.CollectionServ
}

func (r *CollectionRouter) Init(router *gin.RouterGroup, collectionServ collectionserv.CollectionServ) {
	r.collectionServ = collectionServ
	gr := router.Group("collections")
	gr.GET("/", r.GetAllCollections)
	gr.POST("/", r.AddCollection)
	gr.PUT("/", r.UpdateCollection)
	gr.DELETE("/", r.DeleteCollection)
}

// GetAllCollections godoc
// @Summary Get all collections by employee
// @Description Retrieves all collections
// @Tags Collection
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.CollectionResponse
// @Router /employee/collections [get]
func (r *CollectionRouter) GetAllCollections(c *gin.Context) {
	ctx := c.Request.Context()
	cols, err := r.collectionServ.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	colsResp := make([]jsonreqresp.CollectionResponse, len(cols))
	for i, a := range cols {
		colsResp[i] = a.ToCollectionResponse()
	}
	c.JSON(http.StatusOK, colsResp)
}

// AddCollection godoc
// @Summary Add a new collection by employee
// @Description Creates a new collection
// @Tags Collection
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.AddCollectionRequest true "Collection data"
// @Success 201 "Created"
// @Failure 400 "Bad Request"
// @Router /employee/collections [post]
func (r *CollectionRouter) AddCollection(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.AddCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col, err := models.NewCollection(uuid.New(), req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.collectionServ.Add(ctx, &col)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateCollection godoc
// @Summary Update a collection by employee
// @Description Updates an existing collection
// @Tags Collection
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateCollectionRequest true "Collection update data"
// @Success 200  "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Router /employee/collections [put]
func (r *CollectionRouter) UpdateCollection(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.collectionServ.Update(ctx, uuid.MustParse(req.ID), models.CollectionUpdateReq{Title: req.Title})
	if err != nil {
		if errors.Is(err, collectionrep.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteCollection godoc
// @Summary Delete a collection by employee
// @Description Deletes an existing collection
// @Tags Collection
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.DeleteCollectionRequest true "Collection delete data"
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Router /employee/collections [delete]
func (r *CollectionRouter) DeleteCollection(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.DeleteCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.collectionServ.Delete(ctx, uuid.MustParse(req.ID))
	if err != nil {
		if errors.Is(err, collectionrep.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}
