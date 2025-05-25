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
// @Summary Get all events
// @Description Retrieves list of all events
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
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
// @Summary Add new event
// @Description Creates a new event
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body jsonreqresp.AddEventRequest true "Event data"
// @Success 201 "Event created successfully"
// @Failure 400  "Bad Request - Validation error"
// @Failure 401  "unaithorized"
// @Failure 404 "Not Found - Employee not found"
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
// @Summary Delete event
// @Description Deletes existing event
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body jsonreqresp.DeleteEventRequest true "Event delete data"
// @Success 200 "Event deleted successfully"
// @Failure 400 "Bad Request - Validation error"
// @Failure 404 "Not Found - Event not found"
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
// @Summary Update event
// @Description Updates existing event
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body jsonreqresp.UpdateEventRequest true "Event update data"
// @Success 200 "Event updated successfully"
// @Failure 400 "Bad Request - Validation error"
// @Failure 404  "Not Found - Event not found"
// @Router  /employee/events [put]
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
// @Summary Add artwork to event
// @Description Adds an artwork to an existing event
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Event ID"
// @Param request body jsonreqresp.ConArtworkEventRequest true "Artwork to event connection data"
// @Success 200 "Artwork added to event successfully"
// @Failure 400 "Bad Request - Validation error or duplicate artwork"
// @Failure 404 "Not Found - Event or artwork not found"
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
// @Summary Delete artwork from event
// @Description Removes an artwork from an existing event
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Event ID"
// @Param request body jsonreqresp.ConArtworkEventRequest true "Artwork to event connection data"
// @Success 200 "Artwork removed from event successfully"
// @Failure 400 "Bad Request - Validation error"
// @Failure 404 "Not Found - Event or artwork not found"
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
// @Summary Get all artworks from this event by employee
// @Description Retrieves a list of all artworks from this event
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "bearer {token}"
// @Param id path string true "Event ID"
// @Success 200 {array} jsonreqresp.ArtworkResponse
// @Failure 400 "Bad Request - Validation error"
// @Failure 404 "Not Found - Event or artwork not found"
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
