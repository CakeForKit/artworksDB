package frontend

import (
	"net/http"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/components"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend/gintemplrenderer"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/authorserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/collectionserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/eventserv"
	"github.com/gin-gonic/gin"
)

type EmployeeCiteRouter struct {
	authorServ     authorserv.AuthorServ
	collectionServ collectionserv.CollectionServ
	artworkServ    artworkserv.ArtworkService
	eventServ      eventserv.EventService
}

func NewEmployeeCiteRouter(gr *gin.RouterGroup,
	authorServ authorserv.AuthorServ, collectionServ collectionserv.CollectionServ,
	artworkServ artworkserv.ArtworkService, eventServ eventserv.EventService,
) EmployeeCiteRouter {
	r := EmployeeCiteRouter{
		authorServ:     authorServ,
		collectionServ: collectionServ,
		artworkServ:    artworkServ,
		eventServ:      eventServ,
	}
	gr.GET("/authors", r.AuthorsCRUDPage)
	gr.GET("/collections", r.CollectionsCRUDPage)
	gr.GET("/artworks", r.ArtworksCRUDPage)
	gr.GET("/events", r.EventsCRUDPage)

	return r
}

func (r *EmployeeCiteRouter) AuthorsCRUDPage(c *gin.Context) {
	authorsResp := r.allAuthorsResp(c)

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.AuthorsPage(TokenLocalstorage, authorsResp))
	c.Render(http.StatusOK, rend)
}

func (r *EmployeeCiteRouter) CollectionsCRUDPage(c *gin.Context) {
	colsResp := r.allCollectionsResp(c)
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.CollectionsPage(TokenLocalstorage, colsResp))
	c.Render(http.StatusOK, rend)
}

func (r *EmployeeCiteRouter) ArtworksCRUDPage(c *gin.Context) {
	artsResp := r.allArtworksResp(c)
	authorsResp := r.allAuthorsResp(c)
	colsResp := r.allCollectionsResp(c)

	rend := gintemplrenderer.New(
		c.Request.Context(),
		http.StatusOK,
		components.ArtworksCRUDPage(TokenLocalstorage, artsResp, authorsResp, colsResp))
	c.Render(http.StatusOK, rend)
}

func (r *EmployeeCiteRouter) allAuthorsResp(c *gin.Context) []jsonreqresp.AuthorResponse {
	authors, _ := r.authorServ.GetAll(c.Request.Context())
	var authorsResp []jsonreqresp.AuthorResponse
	for _, a := range authors {
		authorsResp = append(authorsResp, a.ToAuthorResponse())
	}
	return authorsResp
}

func (r *EmployeeCiteRouter) allCollectionsResp(c *gin.Context) []jsonreqresp.CollectionResponse {
	cols, _ := r.collectionServ.GetAll(c.Request.Context())
	var colsResp []jsonreqresp.CollectionResponse
	for _, a := range cols {
		colsResp = append(colsResp, a.ToCollectionResponse())
	}
	return colsResp
}

func (r *EmployeeCiteRouter) allArtworksResp(c *gin.Context) []jsonreqresp.ArtworkResponse {
	arts, _ := r.artworkServ.GetAll(c.Request.Context())
	var artsResp []jsonreqresp.ArtworkResponse
	for _, a := range arts {
		artsResp = append(artsResp, a.ToArtworkResponse())
	}
	return artsResp
}

func (r *EmployeeCiteRouter) allEventsResp(c *gin.Context) []jsonreqresp.EventResponse {
	events, _ := r.eventServ.GetAll(c.Request.Context())
	var eventsResp []jsonreqresp.EventResponse
	for _, a := range events {
		eventsResp = append(eventsResp, a.ToEventResponse())
	}
	return eventsResp
}

func (r *EmployeeCiteRouter) EventsCRUDPage(c *gin.Context) {
	artsResp := r.allArtworksResp(c)
	eventsResp := r.allEventsResp(c)

	rend := gintemplrenderer.New(
		c.Request.Context(),
		http.StatusOK,
		components.EventsCRUDPage(TokenLocalstorage, eventsResp, artsResp))
	c.Render(http.StatusOK, rend)
}
