package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth"
	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthUserRouter struct {
	authu auth.AuthUser
}

func (r *AuthUserRouter) Init(router *gin.RouterGroup, authu auth.AuthUser) {
	r.authu = authu
	gr := router.Group("auth-user")
	gr.POST("/register", r.Register)
	gr.POST("/login", r.Login)
	gr.POST("/login-2fa", r.Login2FA)
}

// Register Handler
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags аутентификация
// @Accept json
// @Param request body authmodels.RegisterUserRequest true "Данные для регистрации"
// @Success 200 "Пользователь зарегистрирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Failure 409 "Попытка повторной регистрации"
// @Router /auth-user/register [post]
func (r *AuthUserRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req authmodels.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error %v", err)
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
// @Description Аутентифицирует пользователя, присылает код подтверждения на почту и возвращает id сессии авторизации
// @Tags аутентификация
// @Accept json
// @Param request body authmodels.LoginUserRequest true "Учетные данные для входа"
// @Success 200 "Пользователь успешно аутентифицирован и ему отправлен код подтверждения"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Router /auth-user/login [post]
func (r *AuthUserRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req authmodels.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := r.authu.LoginUser(ctx, req)
	if err != nil {
		if errors.Is(err, userrep.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := authmodels.LoginUserResponse{
		SessionAuthID: sessionID.String(),
	}
	c.JSON(http.StatusOK, rsp)
}

// Login2FA Handler
// @Summary Вход с двухфакторной аутентификацией
// @Description Подтверждает вход пользователя с помощью одноразового кода и выдает access token
// @Tags аутентификация
// @Accept json
// @Param request body authmodels.Login2FAUserRequest true "Данные для подтверждения входа"
// @Success 200 {object} authmodels.Login2FAUserResponse "Успешная аутентификация, возвращает access token"
// @Failure 400 "Неверные входные параметры или невалидный session ID"
// @Failure 401 "Неверный код подтверждения или пользователь не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /auth-user/login-2fa [post]
func (r *AuthUserRouter) Login2FA(c *gin.Context) {
	ctx := c.Request.Context()

	var req authmodels.Login2FAUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(req.SessionAuthID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accessToken, err := r.authu.OTP(ctx, sessionID, req.OTPCode)
	if err != nil {
		if errors.Is(err, userrep.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := authmodels.Login2FAUserResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}
