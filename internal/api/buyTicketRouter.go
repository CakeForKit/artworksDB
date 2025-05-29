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
// @Summary Покупка билетов
// @Description Покупка билетов на указанное мероприятие
// @Tags Билеты
// @Accept json
// @Produce json
// // @Security ApiKeyAuth
// // @Param Authorization header string false "Bearer токен"
// @Param request body jsonreqresp.BuyTicketRequest true "Данные для покупки билетов"
// @Success 200 {object} jsonreqresp.TxTicketPurchaseResponse "Данные покупки сохраняются в cookie"
// @Failure 400 "Неверный формат запроса"
// @Failure 401 "Не авторизован"
// @Failure 404 "Мероприятие не найдено"
// @Failure 409 "Нет доступных билетов"
// @Failure 410 "Транзакция просрочена"
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
// @Summary Получить билеты пользователя
// @Description Получение всех покупок билетов для авторизованного пользователя
// @Tags Билеты
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Success 200 {array} jsonreqresp.TicketPurchaseResponse
// @Failure 401 "Не авторизован"
// @Failure 403 "Доступ запрещен"
// @Router /guest/tickets [get]

// ---
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
// @Summary Подтвердить покупку
// @Description Подтверждает ожидающую транзакцию покупки билетов
// @Tags Билеты
// @Accept json
// @Produce json
// // @Security ApiKeyAuth
// // @Param Authorization header string false "Bearer токен"
// @Param request body jsonreqresp.ConfirmCancelTxRequest true "ID транзакции"
// @Success 200 "Покупка подтверждена"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Транзакция не найдена"
// @Failure 410 "Транзакция просрочена"
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
// @Summary Отменить покупку
// @Description Отменяет ожидающую транзакцию покупки билетов
// @Tags Билеты
// @Accept json
// @Produce json
// // @Security ApiKeyAuth
// // @Param Authorization header string false "Bearer токен"
// @Param request body jsonreqresp.ConfirmCancelTxRequest true "ID транзакции"
// @Success 200 "Покупка отменена"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Транзакция не найдена"
// @Failure 410 "Транзакция просрочена"
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
