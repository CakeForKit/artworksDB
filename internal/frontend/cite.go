package frontend

import (
	"fmt"
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/components"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/gintemplrenderer"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/searcher"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CiteRouter struct {
	searcherServ searcher.Searcher
}

func NewCiteRouter(router *gin.RouterGroup, searcherServ searcher.Searcher) CiteRouter {
	r := CiteRouter{
		searcherServ: searcherServ,
	}
	router.StaticFS("/static", http.Dir("./internal/frontend/static/"))
	gr := router.Group("museum")
	gr.GET("/artworks", r.GetAllArtworks)
	// gr.GET("/artworks/content", r.GetAllArtworksContent)

	return r
}

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
	fmt.Printf("\nsortOps: %+v\n\n", sortOps)
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
		fmt.Printf("%s, ", a.GetTitle())
	}
	fmt.Print("\n")

	fmt.Printf("\nsortOps in: %+v\n\n", sortOps)
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.ArtworksPage(artworksResp, filterOps, sortOps))
	c.Render(http.StatusOK, rend)
}

func (r *CiteRouter) GetAllArtworksContent(c *gin.Context) {
	ctx := c.Request.Context()

	// Читаем параметры из URL (такой же код как в GetAllArtworks)
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

	fmt.Printf("\nsortOps: %+v\n\n", sortOps)
	artworks, err := r.searcherServ.GetAllArtworks(ctx, &filterOps, &sortOps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	artworksResp := make([]jsonreqresp.ArtworkResponse, len(artworks))
	for i, a := range artworks {
		artworksResp[i] = a.ToArtworkResponse()
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.ArtworksContent(artworksResp))
	c.Render(http.StatusOK, rend)
}
