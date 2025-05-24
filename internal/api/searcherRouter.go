package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/searcher"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SearcherRouter struct {
	serv searcher.Searcher
}

func NewSearcherRouter(router *gin.RouterGroup, serv searcher.Searcher) SearcherRouter {
	r := SearcherRouter{
		serv: serv,
	}
	gr := router.Group("museum")
	gr.GET("/artworks", r.GetAllArtworks)
	gr.GET("/events", r.GetAllEvents)
	return r
}

// getArtworks godoc
// @Summary Get artworks
// @Description Retrieves a list of all artworks
// @Tags Searcher
// @Accept json
// @Produce json
// @Param title            query string     false  "Filter by artwork title (max 255 chars)"  maxLength(255)
// @Param author_name      query string     false  "Filter by author name (max 100 chars)"    maxLength(100)
// @Param collection_title query string     false  "Filter by collection title (max 255 chars)" maxLength(255)
// @Param event_id         query string     false  "Filter by event UUID" format(uuid)
// @Param sort_field       query string     true   "Field to sort by"  Enums(title, author_name, creationYear)
// @Param direction_sort   query string     true   "Sort direction"  Enums(ASC, DESC)
// @Success 200 {array} jsonreqresp.ArtworkResponse
// @Router /museum/artworks [get]
func (r *SearcherRouter) GetAllArtworks(c *gin.Context) {
	ctx := c.Request.Context()

	// Читаем параметры из URL
	event_id := uuid.Nil
	if c.Query("event_id") != "" {
		event_id = uuid.MustParse(c.Query("event_id"))
	}
	filterOps := jsonreqresp.ArtworkFilter{
		Title:      c.Query("title"),
		AuthorName: c.Query("author_name"),
		Collection: c.Query("collection_title"),
		EventID:    event_id,
	}

	sortOps := jsonreqresp.ArtworkSortOps{
		Field:     c.Query("sort_field"),
		Direction: c.Query("direction_sort"),
	}

	// Устанавливаем заголовки для кэширования (например, на 1 час)
	// c.Header("Cache-Control", "public, max-age=3600")

	artworks, err := r.serv.GetAllArtworks(ctx, &filterOps, &sortOps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	artworksResp := make([]jsonreqresp.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}
	c.JSON(http.StatusOK, artworksResp)
}

// getAllEvents godoc
// @Summary Get events
// @Description Retrieves a list of all events with optional filtering
// @Tags Searcher
// @Accept json
// @Produce json
// @Param title      query string  false  "Filter by event title"  maxLength(255)
// @Param date_begin query string  false  "Filter by minimum start date (format: YYYY-MM-DD)"  format(date)
// @Param date_end   query string  false  "Filter by maximum end date (format: YYYY-MM-DD)"    format(date)
// @Param can_visit  query boolean false  "Filter by visit availability"
// @Success 200
// @Failure 400 "Invalid date format. Use YYYY-MM-DD"
// @Router /museum/events [get]
func (r *SearcherRouter) GetAllEvents(c *gin.Context) {
	ctx := c.Request.Context()

	parseDate := func(dateStr string) (time.Time, error) {
		return time.Parse("2006-01-02", dateStr)
	}

	filterOps := jsonreqresp.EventFilter{Valid: "true"}
	filterOps.Title = c.Query("title")
	if dateBeginStr := c.Query("date_begin"); dateBeginStr != "" {
		dateBegin, err := parseDate(dateBeginStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_begin format. Use DD-MM-YYYY"})
			return
		}
		filterOps.DateBegin = dateBegin
	}
	if dateEndStr := c.Query("date_end"); dateEndStr != "" {
		dateEnd, err := parseDate(dateEndStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_begin format. Use DD-MM-YYYY"})
			return
		}
		filterOps.DateEnd = dateEnd
	}
	if canVisitStr := c.Query("can_visit"); canVisitStr != "" {
		_, err := strconv.ParseBool(canVisitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid can_visit value (use true/false)"})
			return
		}
		filterOps.CanVisit = canVisitStr
	}

	events, err := r.serv.GetAllEvents(ctx, &filterOps)
	if err != nil {
		if errors.Is(err, jsonreqresp.ErrEventFilterDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	eventsResp := make([]jsonreqresp.EventResponse, len(events))
	for i, a := range events {
		eventsResp[i] = a.ToEventResponse()
	}
	c.JSON(http.StatusOK, eventsResp)
}
