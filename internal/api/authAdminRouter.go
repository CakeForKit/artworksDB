package api

import (
	"errors"
	"net/http"

	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	authadmin "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_admin"
	"github.com/gin-gonic/gin"
)

type AuthAdminRouter struct {
	authu authadmin.AdminAuth
}

func (r *AuthAdminRouter) Init(router *gin.RouterGroup, authu authadmin.AdminAuth) {
	r.authu = authu
	gr := router.Group("auth-admin")
	gr.POST("/register", r.Register)
	gr.POST("/login", r.Login)
}

// Register Handler
// @Summary Register admin
// @Description Register a new admin
// @Tags Аутентификация
// @Accept json
// @Param request body authadmin.RegisterAdminRequest true "Register credentials"
// @Success 200 "The admin registered"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Router /auth-admin/register [post]
func (r *AuthAdminRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req authadmin.RegisterAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.authu.RegisterAdmin(ctx, req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Login Handler
// @Summary Вход администратора
// @Description Аутентифицирует администратора и возвращает токен доступа
// @Tags Аутентификация
// @Accept json
// @Param request body authadmin.LoginAdminRequest true "Учетные данные для входа"
// @Success 200 "Администратор успешно аутентифицирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Failure 403 "Нет прав доступа"
// @Router /auth-admin/login [post]
func (r *AuthAdminRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req authadmin.LoginAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := r.authu.LoginAdmin(ctx, req)
	if err != nil {
		if errors.Is(err, adminrep.ErrAdminNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, authadmin.ErrAdminNotValid) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := authadmin.LoginAdminResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}
