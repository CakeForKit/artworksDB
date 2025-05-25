package frontend

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/components"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/gintemplrenderer"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
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
	// gr.GET("/employee/authors", r.AuthorsCRUDPage)

	return r
}

func (r *CiteRouter) ShowEmployeeLoginPage(c *gin.Context) {
	errorMsg := c.Query("error")
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.LoginPage(token_localstorage, errorMsg))
	c.Render(http.StatusOK, rend)
}

// func (r *CiteRouter) AuthorsCRUDPage(c *gin.Context) {
// 	// Получите список авторов из вашего сервиса
// 	authors, _ := r.authorServ.GetAll(c.Request.Context())

// 	// Преобразуйте в AuthorResponse
// 	var authorsResp []jsonreqresp.AuthorResponse
// 	for _, a := range authors {
// 		authorsResp = append(authorsResp, a.ToAuthorResponse())
// 	}

// 	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.AuthorsPage(token_localstorage, authorsResp))
// 	c.Render(http.StatusOK, rend)
// }

func (r *CiteRouter) GetAllArtworks(c *gin.Context) {
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
		return
	}
	artworksResp := make([]jsonreqresp.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.ArtworksPage(artworksResp, filterOps, sortOps))
	c.Render(http.StatusOK, rend)
}

func (r *CiteRouter) GetAllEvents(c *gin.Context) {
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

	events, err := r.searcherServ.GetAllEvents(ctx, &filterOps)
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

	// fmt.Printf("\nfilterOps: %+v\n\n", filterOps)
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.EventsPage(eventsResp, filterOps))
	c.Render(http.StatusOK, rend)
}
