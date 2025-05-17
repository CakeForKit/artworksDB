package api

import (
	"net/http"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/userservice"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userServ userservice.UserService
}

func NewUserRouter(router *gin.RouterGroup, userServ userservice.UserService) UserRouter {
	r := UserRouter{
		userServ: userServ,
	}

	gr := router.Group("self")
	gr.GET("", r.GetSelf)
	gr.PUT("", r.ChangeSubscribeToMailing)

	return r
}

// GetSelf retrieves current user's profile
// @Summary Get user profile
// @Description Returns authenticated user's profile information
// @Tags User
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {object} jsonreqresp.UserSelfResponse
// @Router /user/self [get]
func (r *UserRouter) GetSelf(c *gin.Context) {
	ctx := c.Request.Context()

	user, err := r.userServ.GetSelf(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := user.ToUserSelfResponse()
	c.JSON(http.StatusOK, resp)
}

// ChangeSubscribeToMailing updates user's mailing subscription preference
// @Summary Update mailing subscription
// @Description Changes user's subscription to email mailings
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param request body jsonreqresp.ChangeSubscribeToMailingRequest true "Subscription preference"
// @Success 200 "OK"
// @Failure 400 "Invalid request body"
// @Router /user/self [put]
func (r *UserRouter) ChangeSubscribeToMailing(c *gin.Context) {
	ctx := c.Request.Context()

	var req jsonreqresp.ChangeSubscribeToMailingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := r.userServ.ChangeSubscribeToMailing(ctx, req.Subscribe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
