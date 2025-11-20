package api

import (
	"errors"
	"net/http"

	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	authemployee "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_employee"
	"github.com/gin-gonic/gin"
)

type AuthEmployeeRouter struct {
	authe authemployee.AuthEmployee
}

func (r *AuthEmployeeRouter) Init(router *gin.RouterGroup, authu authemployee.AuthEmployee) {
	r.authe = authu
	gr := router.Group("auth-employee")

	gr.POST("/login", r.Login)
	gr.POST("/register", r.Register)
}

// Login Handler
// @Summary Вход сотрудника
// @Description Аутентифицирует сотрудника и возвращает токен доступа
// @Tags Аутентификация
// @Accept json
// @Param request body authemployee.LoginEmployeeRequest true "Учетные данные для входа"
// @Success 200 "Сотрудник успешно аутентифицирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Failure 403 "Нет прав доступа"
// @Router /auth-employee/login [post]
func (r *AuthEmployeeRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req authemployee.LoginEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := r.authe.LoginEmployee(ctx, req)
	if err != nil {
		if errors.Is(err, employeerep.ErrEmployeeNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, authemployee.ErrEmployeeNotValid) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := authemployee.LoginEmployeeResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}

// Login Handler
// @Summary Регистрация сотрудника
// @Description Регистрирует сотрудника
// @Tags Аутентификация
// @Accept json
// @Param request body authemployee.RegisterEmployeeRequest true "Учетные данные для регистрации"
// @Success 200 "Сотрудник успешно зарегистрирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Failure 403 "Нет прав доступа"
// @Router /auth-employee/register [post]
func (r *AuthEmployeeRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req authemployee.RegisterEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.authe.RegisterEmployee(ctx, req, req.AdminID)
	if err != nil {
		if errors.Is(err, employeerep.ErrEmployeeNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, authemployee.ErrEmployeeNotValid) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
