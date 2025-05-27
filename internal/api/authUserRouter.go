package api

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/gin-gonic/gin"
)

type AuthUserRouter struct {
	authu auth.AuthUser
}

func (r *AuthUserRouter) Init(router *gin.RouterGroup, authu auth.AuthUser) {
	r.authu = authu
	gr := router.Group("auth-user")
	gr.POST("/register", r.Register)
	gr.POST("/login", r.Login)
}

// Register Handler
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags аутентификация
// @Accept json
// @Param request body auth.RegisterUserRequest true "Данные для регистрации"
// @Success 200 "Пользователь зарегистрирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Failure 409 "Попытка повторной регистрации"
// @Router /auth-user/register [post]

// --------
func (r *AuthUserRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req auth.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.authu.RegisterUser(ctx, req); err != nil {
		if errors.Is(err, userrep.ErrDuplicateLoginUser) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Login Handler
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает токен доступа
// @Tags аутентификация
// @Accept json
// @Param request body auth.LoginUserRequest true "Учетные данные для входа"
// @Success 200 "Пользователь успешно аутентифицирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Router /auth-user/login [post]

// --------
func (r *AuthUserRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req auth.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := r.authu.LoginUser(ctx, req)
	if err != nil {
		if errors.Is(err, userrep.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := auth.LoginUserResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}
