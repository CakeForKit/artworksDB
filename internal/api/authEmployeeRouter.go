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
// @Summary Login employee
// @Description Authenticates a employee and return access token
// @Tags auth-employee
// @Accept json
// @Param request body auth.LoginEmployeeRequest true "Login credentials"
// @Success 200 "The employee has been authenticated"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Failure 403 "Has no rights"
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
