package authemployeeapi

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/gin-gonic/gin"
)

type AuthEmployeeRouter struct {
	authu auth.AuthEmployee
}

func (r *AuthEmployeeRouter) Init(router *gin.RouterGroup, authu auth.AuthEmployee) {
	r.authu = authu
	gr := router.Group("auth-employee")
	gr.POST("/register", r.Register)
	gr.POST("/login", r.Login)
}

// Register Handler
// @Summary Register employee
// @Description Register a new employee
// @Tags auth-employee
// @Accept json
// @Param request body auth.RegisterEmployeeRequest true "Register credentials"
// @Success 200 "The employee registered"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Router /auth-employee/register [post]
func (r *AuthEmployeeRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req auth.RegisterEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.authu.RegisterEmployee(ctx, req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
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
// @Router /auth-employee/login [post]
func (r *AuthEmployeeRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req auth.LoginEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := r.authu.LoginEmployee(ctx, req)
	if err != nil {
		if errors.Is(err, employeerep.ErrEmployeeNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
