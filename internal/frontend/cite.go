package frontend

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/components"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/gintemplrenderer"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/authorserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/searcher"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CiteRouter struct {
	searcherServ searcher.Searcher
	authorServ   authorserv.AuthorServ
}

func NewCiteRouter(router *gin.RouterGroup, searcherServ searcher.Searcher, authorServ authorserv.AuthorServ) CiteRouter {
	r := CiteRouter{
		searcherServ: searcherServ,
		authorServ:   authorServ,
	}

	gr := router.Group("/")
	gr.GET("/artworks", r.GetAllArtworks)
	gr.GET("/events", r.GetAllEvents)
	gr.GET("/login", r.ShowEmployeeLoginPage)
	gr.GET("/events/:id", r.GetEvent)

	return r
}

func (r *CiteRouter) ShowEmployeeLoginPage(c *gin.Context) {
	errorMsg := c.Query("error")
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.LoginPage(TokenLocalstorage, errorMsg))
	c.Render(http.StatusOK, rend)
}

func (r *CiteRouter) allArtworksResp(c *gin.Context) (
	[]jsonreqresp.ArtworkResponse, jsonreqresp.ArtworkFilter, jsonreqresp.ArtworkSortOps,
) {
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

	artworks, err := r.searcherServ.GetAllArtworks(ctx, &filterOps, &sortOps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, jsonreqresp.ArtworkFilter{}, jsonreqresp.ArtworkSortOps{}
	}
	artworksResp := make([]jsonreqresp.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}
	return artworksResp, filterOps, sortOps
}

func (r *CiteRouter) GetAllArtworks(c *gin.Context) {
	artworksResp, filterOps, sortOps := r.allArtworksResp(c)
	if artworksResp != nil {
		rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.ArtworksPage(artworksResp, filterOps, sortOps))
		c.Render(http.StatusOK, rend)
	}
}

// func (r *CiteRouter) GetAllArtworksEmpl(c *gin.Context) {
// 	artworksResp, filterOps, sortOps := r.allArtworksResp(c)
// 	if artworksResp != nil {
// 		rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.EmplArtworksPage(artworksResp, filterOps, sortOps))
// 		c.Render(http.StatusOK, rend)
// 	}
// }

func (r *CiteRouter) allEventsResp(c *gin.Context) (
	[]jsonreqresp.EventResponse, jsonreqresp.EventFilter,
) {
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
			return nil, jsonreqresp.EventFilter{}
		}
		filterOps.DateBegin = dateBegin
	}
	if dateEndStr := c.Query("date_end"); dateEndStr != "" {
		dateEnd, err := parseDate(dateEndStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_begin format. Use DD-MM-YYYY"})
			return nil, jsonreqresp.EventFilter{}
		}
		filterOps.DateEnd = dateEnd
	}
	if canVisitStr := c.Query("can_visit"); canVisitStr != "" {
		_, err := strconv.ParseBool(canVisitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid can_visit value (use true/false)"})
			return nil, jsonreqresp.EventFilter{}
		}
		filterOps.CanVisit = canVisitStr
	}

	events, err := r.searcherServ.GetAllEvents(ctx, &filterOps)
	if err != nil {
		if errors.Is(err, jsonreqresp.ErrEventFilterDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return nil, jsonreqresp.EventFilter{}
	}
	eventsResp := make([]jsonreqresp.EventResponse, len(events))
	for i, a := range events {
		eventsResp[i] = a.ToEventResponse()
	}
	return eventsResp, filterOps
}

func (r *CiteRouter) GetAllEvents(c *gin.Context) {
	eventsResp, filterOps := r.allEventsResp(c)
	if eventsResp != nil {
		rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.EventsPage(eventsResp, filterOps))
		c.Render(http.StatusOK, rend)
	}
}

func (r *CiteRouter) GetEvent(c *gin.Context) {
	ctx := c.Request.Context()
	// Получаем eventID из параметра пути
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}

	event, err := r.searcherServ.GetEvent(ctx, eventID)
	if err != nil {
		if errors.Is(err, eventrep.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	artworks, err := r.searcherServ.GetArtworksFromEvent(ctx, eventID)
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

	rend := gintemplrenderer.New(
		c.Request.Context(),
		http.StatusOK,
		components.EventDetailsPage(
			TokenLocalstorage, event.GetTitle(),
			event.ToEventResponse(), artworksResp,
		),
	)
	c.Render(http.StatusOK, rend)
}

// func (r *CiteRouter) GetAllEventsEmpl(c *gin.Context) {
// 	eventsResp, filterOps := r.allEventsResp(c)
// 	if eventsResp != nil {
// 		rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.EmplEventsPage(eventsResp, filterOps))
// 		c.Render(http.StatusOK, rend)
// 	}
// }
