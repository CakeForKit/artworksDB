package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// _ "github.com/CakeForKit/artworksDB.git/docs"
	"github.com/CakeForKit/artworksDB.git/internal/api"
	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/middleware"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/artworkrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/buyticketstxrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/eventrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/ticketpurchasesrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/adminserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/artworkserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth"
	attemptsrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/attempts_rep"
	authsessionrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_session_repository"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/otp"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/CakeForKit/artworksDB.git/internal/services/authorserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/buyticketserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/collectionserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/eventserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/mailing"
	"github.com/CakeForKit/artworksDB.git/internal/services/searcher"
	"github.com/CakeForKit/artworksDB.git/internal/services/userservice"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getConfig() (
	appCnfg *cnfg.AppConfig,
	dbCreds *cnfg.DatebaseCredentials,
	redisCreds *cnfg.RedisCredentials,
	dbCnfg *cnfg.DatebaseConfig,
	err error,
) {
	appCnfg, err = cnfg.LoadAppConfig("./configs/", "config", "yaml")
	if err != nil {
		return
	}
	if appCnfg.Datebase == cnfg.PostgresDB {
		dbCreds, err = cnfg.LoadPgCredentials("./configs/", "test_db", "env")
		if err != nil {
			err = fmt.Errorf("cannot load PgCredentials: %v", err)
			return
		}
	} else if appCnfg.Datebase == cnfg.ClickHouseDB {
		dbCreds, err = cnfg.LoadClickHouseCredentials("./configs/", "clickhouse", "env")
		if err != nil {
			err = fmt.Errorf("cannot load ClickHouseCredentials: %v", err)
			return
		}
	}

	redisCreds, err = cnfg.LoadRedisCredentials("./configs/", "redis", "env")
	if err != nil {
		panic(fmt.Errorf("cannot load RedisCredentials: %v", err))
	}
	dbCnfg, err = cnfg.LoadDatebaseConfig("./configs/", "config", "yaml")
	if err != nil {
		panic(fmt.Errorf("cannot load DatebaseConfig: %v", err))
	}
	return // appCnfg, dbCreds, redisCreds, dbCnfg
}

func getRepositories(
	ctx context.Context,
	appCnfg *cnfg.AppConfig,
	dbCreds *cnfg.DatebaseCredentials,
	redisCreds *cnfg.RedisCredentials,
	dbCnfg *cnfg.DatebaseConfig,
) (
	userRep userrep.UserRep, employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep, collectionRep collectionrep.CollectionRep,
	authorRep authorrep.AuthorRep, artworkRep artworkrep.ArtworkRep,
	eventRep eventrep.EventRep, txRep buyticketstxrep.BuyTicketsTxRep,
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep,
	err error,
) {
	userRep, err = userrep.NewUserRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	employeeRep, err = employeerep.NewEmployeeRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	adminRep, err = adminrep.NewAdminRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	collectionRep, err = collectionrep.NewCollectionRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	authorRep, err = authorrep.NewAuthorRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	artworkRep, err = artworkrep.NewArtworkRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	eventRep, err = eventrep.NewEventRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	txRep, err = buyticketstxrep.NewBuyTicketsTxRep(ctx, redisCreds)
	if err != nil {
		return
	}
	tPurchasesRep, err = ticketpurchasesrep.NewTicketPurchasesRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	return
}

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
	appCnfg, dbCreds, redisCreds, dbCnfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	// ------------------

	// ----- Repositories -----
	userRep, employeeRep, adminRep, collectionRep,
		authorRep, artworkRep, eventRep, txRep,
		tPurchasesRep, err := getRepositories(ctx, appCnfg, dbCreds, redisCreds, dbCnfg)
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
	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	if err != nil {
		panic(fmt.Errorf("cannot create token maker: %w", err))
	}
	hasher, err := hasher.NewHasher()
	if err != nil {
		panic(err)
	}
	authUserSessionRep := authsessionrep.NewAuthUserSessionRep(appCnfg.DurationLoginSession)
	loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(appCnfg.MaxLoginAttemps, appCnfg.DurationLoginSession)
	otpAttemptRep := attemptsrep.NewOTPAttemptRep(appCnfg.MaxLoginAttemps, appCnfg.DurationLoginSession)
	emailCnfg := cnfg.LoadEmailCnfg()
	emailServ := emailserv.NewEmailService(
		emailCnfg.Host,
		emailCnfg.Port,
		emailCnfg.Username,
		emailCnfg.Password,
		emailCnfg.From,
	)
	otpServ := otp.NewOTPService(*emailServ)
	authUserServ, err := auth.NewAuthUser(
		*appCnfg, userRep,
		tokenMaker, hasher,
		otpServ, authUserSessionRep,
		loginAttemptRep, otpAttemptRep)
	if err != nil {
		panic(err)
	}
	authEmployeeServ, err := auth.NewAuthEmployee(*appCnfg, employeeRep, tokenMaker, hasher)
	if err != nil {
		panic(err)
	}
	authAdminServ, err := auth.NewAuthAdmin(*appCnfg, adminRep, tokenMaker, hasher)
	if err != nil {
		panic(err)
	}
	// serv
	userServ := userservice.NewUserService(userRep, authZ, hasher)
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

	// engine.Run(":8080")
	_ = engine.Run(fmt.Sprintf(":%d", appCnfg.Port))
}
