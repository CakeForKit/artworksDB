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
// @Summary Получить все произведения (сотрудник)
// @Description Возвращает список всех произведений искусства
// @Tags Экспонаты
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.ArtworkResponse
// @Router /employee/artworks [get]
func (r *ArtworksRouter) GetAllArtworks(c *gin.Context) {
	ctx := c.Request.Context()
	artworks, err := r.artworksServ.GetAll(ctx)
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
// @Summary Добавить произведение (сотрудник)
// @Description Добавляет произведение с уже созданными автором и коллекцией
// @Tags Экспонаты
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.AddArtworkRequest true "Новое произведение с существующими автором и коллекцией"
// @Success 201 "Произведение добавлено"
// @Failure 400 "Неверные входные параметры"
// @Failure 404 "Не найдено"
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
// @Summary Обновить произведение (сотрудник)
// @Description Обновляет произведение с новыми/существующими автором и коллекцией
// @Tags Экспонаты
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateArtworkRequest true "Обновляемые данные произведения с автором и коллекцией"
// @Success 200 "Произведение обновлено"
// @Failure 400 "Неверные входные параметры"
// @Failure 404 "Не найдено"
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
// @Summary Удалить произведение (сотрудник)
// @Description Удаляет существующее произведение искусства
// @Tags Экспонаты
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.DeleteArtworkRequest true "Данные для удаления произведения"
// @Success 200 "Успешно удалено"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Не найдено"
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
