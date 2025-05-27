package api

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/authorserv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorRouter struct {
	authorServ authorserv.AuthorServ
}

func NewAuthorRouter(router *gin.RouterGroup, authorServ authorserv.AuthorServ) AuthorRouter {
	r := AuthorRouter{
		authorServ: authorServ,
	}
	gr := router.Group("authors")
	gr.GET("", r.GetAllAuthors)
	gr.POST("", r.AddAuthor)
	gr.PUT("", r.UpdateAuthor)
	gr.DELETE("", r.DeleteAuthor)
	return r
}

// GetAllAuthors godoc
// @Summary Получить всех авторов (сотрудник)
// @Description Возвращает список всех авторов
// @Tags Авторы
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {array} jsonreqresp.AuthorResponse
// @Router /employee/authors [get]
func (r *AuthorRouter) GetAllAuthors(c *gin.Context) {
	ctx := c.Request.Context()
	authors, err := r.authorServ.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	authorsResp := make([]jsonreqresp.AuthorResponse, len(authors))
	for i, a := range authors {
		authorsResp[i] = a.ToAuthorResponse()
	}
	c.JSON(http.StatusOK, authorsResp)
}

// AddAuthor godoc
// @Summary Добавить нового автора (сотрудник)
// @Description Создает нового автора
// @Tags Авторы
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.AddAuthorRequest true "Данные автора"
// @Success 201 "Автор создан"
// @Failure 400 "Неверный запрос"
// @Router /employee/authors [post]
func (r *AuthorRouter) AddAuthor(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.AddAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := models.NewAuthor(uuid.New(), req.Name, req.BirthYear, req.DeathYear)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.authorServ.Add(ctx, &author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateAuthor godoc
// @Summary Обновить автора (сотрудник)
// @Description Обновляет данные существующего автора
// @Tags Авторы
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateAuthorRequest true "Данные для обновления автора"
// @Success 200 "Успешно обновлено"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Автор не найден"
// @Router /employee/authors [put]
func (r *AuthorRouter) UpdateAuthor(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.authorServ.Update(
		ctx, uuid.MustParse(req.ID),
		models.AuthorUpdateReq{
			Name:      req.Name,
			BirthYear: req.BirthYear,
			DeathYear: req.DeathYear,
		})
	if err != nil {
		if errors.Is(err, authorrep.ErrAuthorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteAuthor godoc
// @Summary Удалить автора (сотрудник)
// @Description Удаляет существующего автора
// @Tags Авторы
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.DeleteAuthorRequest true "Данные для удаления автора"
// @Success 200 "Успешно удалено"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Автор не найден"
// @Failure 409 "Конфликт - у автора есть связанные произведения"
// @Router /employee/authors [delete]
func (r *AuthorRouter) DeleteAuthor(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.DeleteAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.authorServ.Delete(ctx, uuid.MustParse(req.ID))
	if err != nil {
		if errors.Is(err, authorrep.ErrAuthorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, authorserv.ErrHasLinkedArtworks) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
