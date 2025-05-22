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
// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Param request body auth.RegisterUserRequest true "Register credentials"
// @Success 200 "The user registered"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Failure 409 "Attempt to re-register"
// @Router /auth-user/register [post]
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
// @Summary Login user
// @Description Authenticates a user and return access token
// @Tags auth
// @Accept json
// @Param request body auth.LoginUserRequest true "Login credentials"
// @Success 200 "The user has been authenticated"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Router /auth-user/login [post]
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
