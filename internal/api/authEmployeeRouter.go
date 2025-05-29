package api

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/gin-gonic/gin"
)

type AuthEmployeeRouter struct {
	authe auth.AuthEmployee
}

func (r *AuthEmployeeRouter) Init(router *gin.RouterGroup, authu auth.AuthEmployee) {
	r.authe = authu
	gr := router.Group("auth-employee")

	gr.POST("/login", r.Login)
}

// Login Handler
// @Summary Вход сотрудника
// @Description Аутентифицирует сотрудника и возвращает токен доступа
// @Tags Аутентификация
// @Accept json
// @Param request body auth.LoginEmployeeRequest true "Учетные данные для входа"
// @Success 200 "Сотрудник успешно аутентифицирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Failure 403 "Нет прав доступа"
// @Router /auth-employee/login [post]
func (r *AuthEmployeeRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req auth.LoginEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := r.authe.LoginEmployee(ctx, req)
	if err != nil {
		if errors.Is(err, employeerep.ErrEmployeeNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, auth.ErrEmployeeNotValid) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := auth.LoginEmployeeResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}
