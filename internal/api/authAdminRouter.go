package api

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/gin-gonic/gin"
)

type AuthAdminRouter struct {
	authu auth.AuthAdmin
}

func (r *AuthAdminRouter) Init(router *gin.RouterGroup, authu auth.AuthAdmin) {
	r.authu = authu
	gr := router.Group("auth-admin")
	// gr.POST("/register", r.Register)
	gr.POST("/login", r.Login)
}

// // Register Handler
// // @Summary Register admin
// // @Description Register a new admin
// // @Tags auth-admin
// // @Accept json
// // @Param request body auth.RegisterAdminRequest true "Register credentials"
// // @Success 200 "The admin registered"
// // @Failure 400 "Wrong input parameters"
// // @Failure 401 "Auth error"
// // @Router /auth-admin/register [post]
// func (r *AuthAdminRouter) Register(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	var req auth.RegisterAdminRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := r.authu.RegisterAdmin(ctx, req); err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{})
// }

// Login Handler
// @Summary Login admin
// @Description Authenticates a admin and return access token
// @Tags auth
// @Accept json
// @Param request body auth.LoginAdminRequest true "Login credentials"
// @Success 200 "The admin has been authenticated"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Failure 403 "Has no rights"
// @Router /auth-admin/login [post]
func (r *AuthAdminRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req auth.LoginAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := r.authu.LoginAdmin(ctx, req)
	if err != nil {
		if errors.Is(err, adminrep.ErrAdminNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, auth.ErrAdminNotValid) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := auth.LoginAdminResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}
