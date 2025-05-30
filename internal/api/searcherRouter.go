package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
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
	gr.GET("/events/:id", r.GetEvent)
	gr.GET("/events/:id/artworks", r.GetArtworkFromEvent)
	gr.GET("/events/:id/statcols", r.GetCollectionsStat)
	return r
}

// GetEvent godoc
// @Summary Получить мероприятие по ID
// @Description Возвращает одно мероприятие по его идентификатору
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id path string true "ID мероприятия"
// @Success 200 {object} jsonreqresp.EventResponse
// @Failure 400 "Неверный формат ID"
// @Failure 404 "Мероприятие не найдено"
// @Router /museum/events/{id} [GET]
func (r *SearcherRouter) GetEvent(c *gin.Context) {
	ctx := c.Request.Context()
	// Получаем eventID из параметра пути
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}
	event, err := r.serv.GetEvent(ctx, eventID)
	if err != nil {
		if errors.Is(err, eventrep.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, event.ToEventResponse())

}

// GetArtworkFromEvent godoc
// @Summary Получить все произведения мероприятия
// @Description Возвращает список всех произведений данного мероприятия
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id path string true "ID мероприятия"
// @Success 200 {array} jsonreqresp.ArtworkResponse
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 404 "Не найдено - мероприятие или произведение не найдено"
// @Router /museum/events/{id}/artworks [get]
func (r *SearcherRouter) GetArtworkFromEvent(c *gin.Context) {
	ctx := c.Request.Context()
	// Получаем eventID из параметра пути
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}

	artworks, err := r.serv.GetArtworksFromEvent(ctx, eventID)
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

// getArtworks godoc
// @Summary Получить произведения
// @Description Возвращает список всех произведений с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param title            query string     false  "Фильтр по названию произведения (макс. 255 символов)"  maxLength(255)
// @Param author_name      query string     false  "Фильтр по имени автора (макс. 100 символов)"    maxLength(100)
// @Param collection_title query string     false  "Фильтр по названию коллекции (макс. 255 символов)" maxLength(255)
// @Param event_id         query string     false  "Фильтр по ID мероприятия" format(uuid)
// @Param sort_field       query string     true   "Поле для сортировки"  Enums(title, author_name, creationYear)
// @Param direction_sort   query string     true   "Направление сортировки"  Enums(ASC, DESC)
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
// @Summary Получить мероприятия
// @Description Возвращает список всех мероприятий с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param title      query string  false  "Фильтр по названию мероприятия"  maxLength(255)
// @Param date_begin query string  false  "Фильтр по минимальной дате начала (формат: ГГГГ-ММ-ДД)"  format(date)
// @Param date_end   query string  false  "Фильтр по максимальной дате окончания (формат: ГГГГ-ММ-ДД)"    format(date)
// @Param can_visit  query boolean false  "Фильтр по доступности для посещения"
// @Success 200 {array} jsonreqresp.EventResponse
// @Failure 400 "Неверный формат даты. Используйте ГГГГ-ММ-ДД"
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

// GetCollectionsStat godoc
// @Summary Получить статистику по коллекциям для мероприятия
// @Description Возвращает список коллекиций произведения искусства из которых участвуют в выставке
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id path string true "ID мероприятия"
// @Success 200 {array} jsonreqresp.StatCollectionsResponse
// @Failure 400 "Неверный запрос - ошибка валидации"
// @Failure 404 "Не найдено - мероприятие или произведение не найдено"
// @Router /museum/events/{id}/statcols [get]
func (r *SearcherRouter) GetCollectionsStat(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем eventID из параметра пути
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}

	statCols, err := r.serv.GetCollectionsStat(ctx, eventID)
	if err != nil {
		if errors.Is(err, eventrep.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	resp := make([]jsonreqresp.StatCollectionsResponse, len(statCols))
	for i, v := range statCols {
		resp[i] = v.ToResponse()
	}

	c.JSON(http.StatusOK, resp)
}
