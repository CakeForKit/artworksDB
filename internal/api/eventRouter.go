package api

import (
	"errors"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/eventserv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventRouter struct {
	eventServ eventserv.EventService
	authZ     auth.AuthZ
}

func NewEventRouter(router *gin.RouterGroup, eventServ eventserv.EventService, authZ auth.AuthZ) EventRouter {
	r := EventRouter{
		eventServ: eventServ,
		authZ:     authZ,
	}
	gr := router.Group("events")
	gr.GET("/", r.GetAllEvents)
	gr.POST("", r.AddEvent)
	gr.DELETE("", r.DeleteEvent)
	gr.PUT("", r.UpdateEvent)
	gr.PUT("/:id", r.AddArtworkToEvent)
	gr.DELETE("/:id", r.DeleteArtworkFromEvent)
	gr.GET("/:id/artworks", r.GetArtworkFromEvent)
	return r
}

// GetAllEvents godoc
// @Summary Получить все мероприятия (сотрудник)
// @Description Возвращает список всех мероприятий
// @Tags Мероприятия
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Success 200 {array} jsonreqresp.EventResponse
// @Router /employee/events [get]
func (r *EventRouter) GetAllEvents(c *gin.Context) {
	ctx := c.Request.Context()
	events, err := r.eventServ.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eventsResp := make([]jsonreqresp.EventResponse, len(events))
	for i, e := range events {
		eventsResp[i] = e.ToEventResponse()
	}
	c.JSON(http.StatusOK, eventsResp)

}

// AddEvent godoc
// @Summary Добавить новое мероприятие (сотрудник)
// @Description Создает новое мероприятие
// @Tags Мероприятия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body jsonreqresp.AddEventRequest true "Данные мероприятия"
// @Success 201 "Мероприятие успешно создано"
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 401 "Не авторизован"
// @Failure 404 "Не найдено - сотрудник не найден"
// @Router /employee/events [post]
func (r *EventRouter) AddEvent(c *gin.Context) {
	ctx := c.Request.Context()
	var req jsonreqresp.AddEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employeeID, err := r.authZ.EmployeeIDFromContext(ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	addReq := jsonreqresp.EventAdd{
		Title:      req.Title,
		DateBegin:  req.DateBegin,
		DateEnd:    req.DateEnd,
		Address:    req.Address,
		CanVisit:   *req.CanVisit,
		EmployeeID: employeeID,
		CntTickets: *req.CntTickets,
		ArtworkIDs: req.ArtworkIDs,
	}
	if err := r.eventServ.Add(ctx, &addReq); err != nil {
		if errors.Is(err, eventrep.ErrAddNoEmployee) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, models.ErrValidateEvent) || errors.Is(err, eventserv.ErrArtworkBusy) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteEvent godoc
// @Summary Удалить мероприятие (сотрудник)
// @Description Удаляет существующее мероприятие
// @Tags Мероприятия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body jsonreqresp.DeleteEventRequest true "Данные для удаления мероприятия"
// @Success 200 "Мероприятие успешно удалено"
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 404 "Не найдено - мероприятие не найдено"
// @Router /employee/events [delete]
func (r *EventRouter) DeleteEvent(c *gin.Context) {
	ctx := c.Request.Context()
	var req jsonreqresp.DeleteEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.eventServ.Delete(ctx, uuid.MustParse(req.ID)); err != nil {
		if errors.Is(err, eventrep.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateEvent godoc
// @Summary Обновить мероприятие (сотрудник)
// @Description Обновляет существующее мероприятие
// @Tags Мероприятия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body jsonreqresp.UpdateEventRequest true "Данные для обновления мероприятия"
// @Success 200 "Мероприятие успешно обновлено"
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 404 "Не найдено - мероприятие не найдено"
// @Router /employee/events [put]
func (r *EventRouter) UpdateEvent(c *gin.Context) {
	ctx := c.Request.Context()
	var req jsonreqresp.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.eventServ.Update(
		ctx, uuid.MustParse(req.ID),
		&jsonreqresp.EventUpdate{
			Title:      req.Title,
			DateBegin:  req.DateBegin,
			DateEnd:    req.DateEnd,
			Address:    req.Address,
			CanVisit:   *req.CanVisit,
			CntTickets: *req.CntTickets,
			// Valid:      req.Valid,
		})
	if err != nil {
		if errors.Is(err, eventrep.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// AddArtworkToEvent godoc
// @Summary Добавить произведение к мероприятию (сотрудник)
// @Description Добавляет произведение к существующему мероприятию
// @Tags Мероприятия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param id path string true "ID мероприятия"
// @Param request body jsonreqresp.ConArtworkEventRequest true "Данные для связи произведения с мероприятием"
// @Success 200 "Произведение успешно добавлено к мероприятию"
// @Failure 400 "Неверный запрос - ошибка валидации или дублирование произведения"
// @Failure 404 "Не найдено - мероприятие или произведение не найдено"
// @Router /employee/events/{id} [PUT]
func (r *EventRouter) AddArtworkToEvent(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем eventID из параметра пути
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}

	var req jsonreqresp.ConArtworkEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	artworkIDs := uuid.UUIDs{uuid.MustParse(req.ArtworkID)}
	err = r.eventServ.AddArtworksToEvent(ctx, eventID, artworkIDs)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateArtwokIDs) || errors.Is(err, models.ErrAddArtwork) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if errors.Is(err, eventrep.ErrEventNotFound) || errors.Is(err, artworkrep.ErrArtworkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteArtworkFromEvent godoc
// @Summary Удалить произведение из мероприятия (сотрудник)
// @Description Удаляет произведение из существующего мероприятия
// @Tags Мероприятия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param id path string true "ID мероприятия"
// @Param request body jsonreqresp.ConArtworkEventRequest true "Данные для связи произведения с мероприятием"
// @Success 200 "Произведение успешно удалено из мероприятия"
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 404 "Не найдено - мероприятие или произведение не найдено"
// @Router /employee/events/{id} [delete]
func (r *EventRouter) DeleteArtworkFromEvent(c *gin.Context) {
	ctx := c.Request.Context()
	// Получаем eventID из параметра пути
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}

	var req jsonreqresp.ConArtworkEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.eventServ.DeleteArtworkFromEvent(ctx, eventID, uuid.MustParse(req.ArtworkID))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateArtwokIDs) || errors.Is(err, models.ErrAddArtwork) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if errors.Is(err, eventrep.ErrEventNotFound) ||
			errors.Is(err, artworkrep.ErrArtworkNotFound) ||
			errors.Is(err, eventrep.ErrEventArtowrkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetArtworkFromEvent godoc
// @Summary Получить все произведения мероприятия (сотрудник)
// @Description Возвращает список всех произведений данного мероприятия
// @Tags Мероприятия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param id path string true "ID мероприятия"
// @Success 200 {array} jsonreqresp.ArtworkResponse
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 404 "Не найдено - мероприятие или произведение не найдено"
// @Router /employee/events/{id}/artworks [get]
func (r *EventRouter) GetArtworkFromEvent(c *gin.Context) {
	ctx := c.Request.Context()
	// Получаем eventID из параметра пути
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}

	artworks, err := r.eventServ.GetArtworksFromEvent(ctx, eventID)
	if err != nil {
		if errors.Is(err, eventrep.ErrEventNotFound) ||
			errors.Is(err, artworkrep.ErrArtworkNotFound) ||
			errors.Is(err, eventrep.ErrEventArtowrkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	artworksResp := make([]jsonreqresp.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}
	c.JSON(http.StatusOK, artworksResp)
}
