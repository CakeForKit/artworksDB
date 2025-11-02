package api

import (
	"errors"
	"net/http"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/collectionserv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CollectionRouter struct {
	collectionServ collectionserv.CollectionServ
}

func (r *CollectionRouter) Init(router *gin.RouterGroup, collectionServ collectionserv.CollectionServ) {
	r.collectionServ = collectionServ
	gr := router.Group("collections")
	gr.GET("/", r.GetAll)
	gr.POST("/", r.Add)
	gr.PUT("/", r.Update)
	gr.DELETE("/", r.Delete)
}

// GetAll godoc
// @Summary Получить все коллекции (сотрудник)
// @Description Возвращает список всех коллекций
// @Tags Коллекции
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.CollectionResponse
// @Router /employee/collections [get]
func (r *CollectionRouter) GetAll(c *gin.Context) {
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

// Add godoc
// @Summary Добавить новую коллекцию (сотрудник)
// @Description Создает новую коллекцию
// @Tags Коллекции
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.AddRequest true "Данные коллекции"
// @Success 201 "Коллекция создана"
// @Failure 400 "Неверный запрос"
// @Router /employee/collections [post]
func (r *CollectionRouter) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.AddRequest
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

// Update godoc
// @Summary Обновить коллекцию (сотрудник)
// @Description Обновляет существующую коллекцию
// @Tags Коллекции
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateRequest true "Данные для обновления коллекции"
// @Success 200 "Коллекция обновлена"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Коллекция не найдена"
// @Router /employee/collections [put]
func (r *CollectionRouter) Update(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.UpdateRequest
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

// Delete godoc
// @Summary Удалить коллекцию (сотрудник)
// @Description Удаляет существующую коллекцию
// @Tags Коллекции
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.DeleteRequest true "Данные для удаления коллекции"
// @Success 200 "Коллекция удалена"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Коллекция не найдена"
// @Router /employee/collections [delete]
func (r *CollectionRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.DeleteRequest
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
