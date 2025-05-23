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
	gr.GET("/", r.GetAllAuthors)
	gr.POST("/", r.AddAuthor)
	gr.PUT("/", r.UpdateAuthor)
	gr.DELETE("/", r.DeleteAuthor)
	return r
}

// GetAllAuthors godoc
// @Summary Get all authors by employee
// @Description Retrieves all authors
// @Tags Author
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
// @Summary Add a new author by employee
// @Description Creates a new author
// @Tags Author
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.AddAuthorRequest true "Author data"
// @Success 201 "Created"
// @Failure 400 "Bad Request"
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
// @Summary Update an author by employee
// @Description Updates an existing author
// @Tags Author
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.UpdateAuthorRequest true "Author update data"
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
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
// @Summary Delete an author by employee
// @Description Deletes an existing author
// @Tags Author
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.DeleteAuthorRequest true "Author delete data"
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 409 "Conflict - Author has linked artworks"
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
