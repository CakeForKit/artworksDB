package api

import (
	"errors"
	"net/http"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/eventserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/mailing"
	"github.com/gin-gonic/gin"
)

type MailingRouter struct {
	mailingServ mailing.MailingService
	eventServ   eventserv.EventService
}

func NewMailingRouter(router *gin.RouterGroup, mailingServ mailing.MailingService, eventServ eventserv.EventService) MailingRouter {
	r := MailingRouter{
		mailingServ: mailingServ,
		eventServ:   eventServ,
	}
	gr := router.Group("mailing")
	gr.POST("/", r.SendMails)

	return r
}

// SendMails sends a mailing to all users based on events.
// @Summary Send mailing to users
// @Description Sends a message to all users using event data
// @Tags Mailing
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Success 200 {object} jsonreqresp.MailingResponse "Mailing sent successfully"
// @Failure 404 "Error: no events found"
// @Router /employee/mailing/ [post]

// ---
func (r *MailingRouter) SendMails(c *gin.Context) {
	ctx := c.Request.Context()

	events, err := r.eventServ.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": errors.New("no events")})
		return
	}

	msgText, userIDs, err := r.mailingServ.SendMailToAllUsers(ctx, events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := jsonreqresp.MailingResponse{
		MsgText: msgText,
		UserIDs: userIDs.Strings(),
	}
	c.JSON(http.StatusOK, resp)
}
