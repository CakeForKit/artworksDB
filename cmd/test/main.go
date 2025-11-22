package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// _ "github.com/CakeForKit/artworksDB.git/docs"
	"github.com/CakeForKit/artworksDB.git/internal/api"
	"github.com/CakeForKit/artworksDB.git/internal/creators"
	"github.com/CakeForKit/artworksDB.git/internal/middleware"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("TEST MAIN\n\n")
	ctx := context.Background()
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Можно указать конкретные домены вместо "*"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.OPTIONS("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})
	apiGroup := engine.Group("/api/v1")

	// timingTracer, err := tracing.NewTimingTracer(tracing.DefaultTimingConfig())
	// if err != nil {
	// 	panic(err)
	// }
	// defer timingTracer.Shutdown(ctx)
	// apiGroup.Use(middleware.MetricsMiddleware(timingTracer))

	tracer, err := tracing.NewTracer(tracing.DefaultConfigTracer())
	if err != nil {
		panic(err)
	}
	tracer.ChangeEnabled(false)
	// defer tracer.Shutdown(ctx)
	// if tracer.IsEnabled() {
	// 	apiGroup.Use(middleware.TraceMiddleware(tracer))
	// }
	/*
		logCnfg, err := cnfg.GetLogConfig()
		if err != nil {
			panic(fmt.Errorf("cannot load LogConfig: %v", err))
		}
		projLogger, err := projlog.NewLogger(logCnfg)
		if err != nil {
			panic(err)
		}
		defer projLogger.Sync()
		apiGroup.Use(middleware.LogMiddleware(projLogger))
	*/

	// ----- Config ------
	appCnfg, dbCreds, redisCreds, dbCnfg, err := creators.GetConfig()
	if err != nil {
		panic(err)
	}
	// ------------------

	// ----- Repositories -----
	userRep, employeeRep, adminRep, collectionRep,
		authorRep, artworkRep, eventRep, txRep,
		tPurchasesRep, err := creators.GetRepositories(ctx, appCnfg, dbCreds, redisCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	// ------------------------

	// ----- Services -----
	// auth
	authZ, hasher, authUserServ, authEmployeeServ, authAdminServ := creators.GetAuthServs(
		appCnfg, tracer, userRep, employeeRep, adminRep)

	// serv
	userServ, adminServ, buyTicketServ,
		collectionServ, authroServ, artworkServ,
		eventServ, searcherServ, mailingServ := creators.GetServs(
		appCnfg, tracer, userRep, employeeRep, adminRep, authZ, hasher,
		txRep, tPurchasesRep, collectionRep, authorRep, artworkRep, eventRep)

	// --------------------

	// ----- Groups -----

	userGroup := apiGroup.Group("/user")
	userGroup.Use(middleware.AuthMiddleware(authUserServ, authZ, true))
	guestGroup := apiGroup.Group("/guest")
	guestGroup.Use(middleware.AuthMiddleware(authUserServ, authZ, false))
	employeeGroup := apiGroup.Group("/employee")
	employeeGroup.Use(middleware.AuthMiddleware(authEmployeeServ, authZ, true))
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(authAdminServ, authZ, true))
	// ------------------------

	// ----- Routers -----
	authUserRouter := api.AuthUserRouter{}
	authUserRouter.Init(apiGroup, authUserServ)
	authEmployeeRouter := api.AuthEmployeeRouter{}
	authEmployeeRouter.Init(apiGroup, authEmployeeServ)
	authAdminRouter := api.AuthAdminRouter{}
	authAdminRouter.Init(apiGroup, authAdminServ)

	userRouter := api.NewUserRouter(userGroup, userServ)
	_ = userRouter
	employeeRouter := api.AdminRouter{}
	employeeRouter.Init(adminGroup, adminServ, authEmployeeServ, authZ)

	collectionRouter := api.CollectionRouter{}
	collectionRouter.Init(employeeGroup, collectionServ)
	authorRouter := api.NewAuthorRouter(employeeGroup, authroServ)
	_ = authorRouter
	artworkRouter := api.NewArtworksRouter(employeeGroup, artworkServ)
	_ = artworkRouter
	eventRouter := api.NewEventRouter(employeeGroup, eventServ, authZ)
	_ = eventRouter
	mailingRouter := api.NewMailingRouter(employeeGroup, mailingServ, eventServ)
	_ = mailingRouter
	buyTicketRouter := api.NewBuyTicketRouter(guestGroup, buyTicketServ)
	_ = buyTicketRouter
	searcherRouter := api.NewSearcherRouter(apiGroup, searcherServ)
	_ = searcherRouter
	// -------------------

	// engine.Run(":8080")
	_ = engine.Run(fmt.Sprintf(":%d", appCnfg.Port))
}
