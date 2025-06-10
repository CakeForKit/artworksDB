// @title Музей
// @version 1.0
// @description API для системы учета произведений искусств
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "git.iu7.bmstu.ru/ped22u691/PPO.git/docs"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/api"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/frontend"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/middleware"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/buyticketstxrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/ticketpurchasesrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/adminserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/authorserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/buyticketserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/collectionserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/eventserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/mailing"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/searcher"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/userservice"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main2() {
	ctx := context.Background()
	clhCreds, err := cnfg.LoadClickHouseCredentials()
	if err != nil {
		panic(fmt.Errorf("cannot load ClickHouseCredentials: %v", err))
	}
	dbCnfg, err := cnfg.LoadDatebaseConfig("./configs/")
	if err != nil {
		panic(fmt.Errorf("cannot load DatebaseConfig: %v", err))
	}
	appCnfg, err := cnfg.LoadAppConfig()
	if err != nil {
		panic(fmt.Errorf("cannot load AppConfig: %v", err))
	}
	arep, err := eventrep.NewEventRep(ctx, appCnfg.Datebase, clhCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	res, err := arep.GetAll(ctx, &jsonreqresp.EventFilter{})
	if err != nil {
		panic(err)
	}
	for _, a := range res {
		fmt.Printf("%+v\n\n", *a)
	}

}

func main() {
	ctx := context.Background()
	engine := gin.New()
	// engine.Handle("OPTIONS", "/api/*path", func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "*")
	// 	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Content-Type")
	// 	c.Status(200)
	// })
	// Настройка CORS
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

	// logCnfg, err := cnfg.GetLogConfig()
	// if err != nil {
	// 	panic(fmt.Errorf("cannot load LogConfig: %v", err))
	// }
	// projLogger, err := projlog.NewLogger(logCnfg)
	// if err != nil {
	// 	panic(err)
	// }
	// defer projLogger.Sync()
	// apiGroup.Use(middleware.LogMiddleware(projLogger))

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// ----- Config ------
	appCnfg, err := cnfg.LoadAppConfig()
	if err != nil {
		panic(fmt.Errorf("cannot load AppConfig: %v", err))
	}
	var dbCreds *cnfg.DatebaseCredentials
	if appCnfg.Datebase == cnfg.PostgresDB {
		pgCreds, err := cnfg.LoadPgCredentials("./configs/")
		if err != nil {
			panic(fmt.Errorf("cannot load PgCredentials: %v", err))
		}
		dbCreds = pgCreds
	} else if appCnfg.Datebase == cnfg.ClickHouseDB {
		clhCreds, err := cnfg.LoadClickHouseCredentials()
		if err != nil {
			panic(fmt.Errorf("cannot load ClickHouseCredentials: %v", err))
		}
		dbCreds = clhCreds
	}
	redisCreds, err := cnfg.LoadRedisCredentials()
	if err != nil {
		panic(fmt.Errorf("cannot load RedisCredentials: %v", err))
	}
	dbCnfg, err := cnfg.LoadDatebaseConfig("./configs/")
	if err != nil {
		panic(fmt.Errorf("cannot load DatebaseConfig: %v", err))
	}
	// ------------------

	// для Swagger - НЕ ТРОГАТЬ
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", appCnfg.Port))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// ----- Repositories -----
	userRep, err := userrep.NewUserRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	employeeRep, err := employeerep.NewEmployeeRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	adminRep, err := adminrep.NewAdminRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	collectionRep, err := collectionrep.NewCollectionRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	authorRep, err := authorrep.NewAuthorRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	artworkRep, err := artworkrep.NewArtworkRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	eventRep, err := eventrep.NewEventRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	txRep, err := buyticketstxrep.NewBuyTicketsTxRep(ctx, redisCreds)
	if err != nil {
		panic(err)
	}
	tPurchasesRep, err := ticketpurchasesrep.NewTicketPurchasesRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	// ------------------------

	// ----- Services -----
	// auth
	authZ, err := auth.NewAuthZ()
	if err != nil {
		panic(err)
	}
	authUserServ, err := auth.NewAuthUser(*appCnfg, userRep)
	if err != nil {
		panic(err)
	}
	authEmployeeServ, err := auth.NewAuthEmployee(*appCnfg, employeeRep)
	if err != nil {
		panic(err)
	}
	authAdminServ, err := auth.NewAuthAdmin(*appCnfg, adminRep)
	if err != nil {
		panic(err)
	}
	// serv
	userServ := userservice.NewUserService(userRep, authZ)
	adminserv := adminserv.NewAdminService(employeeRep, userRep, authZ)
	buyTicketServ, _ := buyticketserv.NewBuyTicketsServ(txRep, tPurchasesRep, *appCnfg, authZ, userRep, eventRep)
	collectionServ := collectionserv.NewCollectionServ(collectionRep)
	authroServ := authorserv.NewAuthorServ(authorRep)
	artworkServ := artworkserv.NewArtworkService(artworkRep, authorRep, collectionRep)
	eventServ := eventserv.NewEventService(eventRep, artworkRep)
	searcherServ := searcher.NewSearcher(artworkRep, eventRep)
	mailingServ := mailing.NewGmailSender(userRep, "museum", "museum@test.ru", "1234")
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
	employeeRouter.Init(adminGroup, adminserv, authEmployeeServ, authZ)

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

	// ------ Cite -----
	engine.StaticFS("/static", http.Dir("./internal/frontend/static/"))
	citeGroup := engine.Group("museum")
	citeRouter := frontend.NewCiteRouter(citeGroup, searcherServ, authroServ)
	_ = citeRouter
	emplCiteGroup := citeGroup.Group("employee")
	emplCiteGroup.Use(middleware.AuthMiddleware(authEmployeeServ, authZ, true))
	employeesCiteRouter := frontend.NewEmployeeCiteRouter(
		emplCiteGroup, authroServ, collectionServ, artworkServ, eventServ)
	_ = employeesCiteRouter

	// // Статические файлы
	// citeGroup.Static("/static", filepath.Join("internal", "frontend", "static"))
	// ---------------

	// engine.Run(":8080")
	engine.Run(fmt.Sprintf(":%d", appCnfg.Port))
}
