package api

import (
	"encoding/json"
	"errors"
	"net/http"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/buyticketserv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	DataTicketPurchaseTx = "DataTicketPurchaseTx"
)

type BuyTicketRouter struct {
	buyTicketServ buyticketserv.BuyTicketsServ
}

func NewBuyTicketRouter(router *gin.RouterGroup, buyTicketServ buyticketserv.BuyTicketsServ) BuyTicketRouter {
	r := BuyTicketRouter{
		buyTicketServ: buyTicketServ,
	}
	gr := router.Group("tickets")
	gr.POST("", r.BuyTickets)
	gr.GET("", r.GetAllTicketPurchasesOfUser)
	gr.PUT("/confirm", r.ConfirmBuyTicket)
	gr.PUT("/cancel", r.CancelBuyTicket)
	return r
}

// BuyTickets purchases tickets for an event
// @Summary Purchase tickets
// @Description Buy tickets for a specific event
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "Bearer token"
// @Param request body jsonreqresp.BuyTicketRequest true "Ticket purchase details"
// @Success 200 {object} jsonreqresp.TxTicketPurchaseResponse "Sets purchase data in cookie"
// @Failure 400 "Invalid request format"
// @Failure 401 "Unauthorized"
// @Failure 404 "Event not found"
// @Failure 409 "No tickets available"
// @Failure 410 "Transaction expired"
// @Router /guest/tickets [post]
func (r *BuyTicketRouter) BuyTickets(c *gin.Context) {
	ctx := c.Request.Context()
	var req jsonreqresp.BuyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txPurchase, err := r.buyTicketServ.BuyTicket(
		ctx, uuid.MustParse(req.EventID), req.CntTickets,
		req.CustomerName, req.CustomerEmail)
	if err != nil {
		if errors.Is(err, buyticketserv.ErrNoFreeTicket) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if errors.Is(err, buyticketserv.ErrNoUserData) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if errors.Is(err, buyticketstxrep.ErrExpireTx) {
			c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	txResp := txPurchase.ToTxTicketPurchaseResponse()
	txData, err := json.Marshal(txResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize purchase data"})
		return
	}
	c.SetCookie(
		DataTicketPurchaseTx, string(txData),
		int(r.buyTicketServ.GetBuyTicketTransactionDuration().Seconds()),
		"/", "", false, true)
	c.JSON(http.StatusOK, txResp)
}

// GetAllTicketPurchasesOfUser retrieves user's ticket purchases
// @Summary Get user's tickets
// @Description Retrieves all ticket purchases for authenticated user
// @Tags Tickets
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} jsonreqresp.TicketPurchaseResponse
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /guest/tickets [get]
func (r *BuyTicketRouter) GetAllTicketPurchasesOfUser(c *gin.Context) {
	ctx := c.Request.Context()

	txPurchases, err := r.buyTicketServ.GetAllTicketPurchasesOfUser(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNotAuthZ) || errors.Is(err, auth.ErrHasNoRights) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	txPurchasesResp := make([]jsonreqresp.TicketPurchaseResponse, len(txPurchases))
	for i, t := range txPurchases {
		txPurchasesResp[i] = t.ToTicketPurchaseResponse()
	}
	c.JSON(http.StatusOK, txPurchasesResp)
}

// ConfirmBuyTicket confirms a ticket purchase
// @Summary Confirm purchase
// @Description Confirms a pending ticket purchase
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "Bearer token"
// @Param request body jsonreqresp.ConfirmCancelTxRequest true "Transaction ID"
// @Success 200 "Purchase confirmed"
// @Failure 400 "Invalid request"
// @Failure 404 "Transaction not found"
// @Failure 410 "Transaction expired"
// @Router /guest/tickets/confirm [put]
func (r *BuyTicketRouter) ConfirmBuyTicket(c *gin.Context) {
	ctx := c.Request.Context()
	var req jsonreqresp.ConfirmCancelTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.buyTicketServ.ConfirmBuyTicket(ctx, uuid.MustParse(req.TxID))
	if err != nil {
		if errors.Is(err, buyticketstxrep.ErrTxNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, buyticketstxrep.ErrExpireTx) {
			c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// CancelBuyTicket cancels a ticket purchase
// @Summary Cancel purchase
// @Description Cancels a pending ticket purchase
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "Bearer token"
// @Param request body jsonreqresp.ConfirmCancelTxRequest true "Transaction ID"
// @Success 200 "Purchase cancelled"
// @Failure 400 "Invalid request"
// @Failure 404 "Transaction not found"
// @Failure 410 "Transaction expired"
// @Router /guest/tickets/cancel [put]
func (r *BuyTicketRouter) CancelBuyTicket(c *gin.Context) {
	ctx := c.Request.Context()
	var req jsonreqresp.ConfirmCancelTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.buyTicketServ.CancelBuyTicket(ctx, uuid.MustParse(req.TxID))
	if err != nil {
		if errors.Is(err, buyticketstxrep.ErrTxNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, buyticketstxrep.ErrExpireTx) {
			c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
