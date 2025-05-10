// @title Museum
// @version 1.0
// @description API для системы учета произведений искусств
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	_ "git.iu7.bmstu.ru/ped22u691/PPO.git/docs"
	artworksapi "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/api/artworks_api"
	authadminapi "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/api/auth_admin_api"
	authemployeeapi "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/api/auth_employee_api"
	authuserapi "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/api/auth_user_api"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/middleware"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/artworkserv"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	ctx := context.Background()
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	apiGroup := engine.Group("/api/v1")

	// для Swagger - НЕ ТРОГАТЬ
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// ----- Config ------
	pgCreds, err := cnfg.LoadPgCredentials()
	if err != nil {
		panic(fmt.Errorf("cannot load PgCredentials: %v", err))
	}
	dbCnfg, err := cnfg.LoadDatebaseConfig("./configs/")
	if err != nil {
		panic(fmt.Errorf("cannot load DatebaseConfig: %v", err))
	}
	appCnfg, err := cnfg.LoadAppConfig()
	if err != nil {
		panic(fmt.Errorf("cannot load AppConfig: %v", err))
	}
	// ------------------

	// User
	{
		// ----- User Auth -----
		userRep, err := userrep.NewUserRep(ctx, pgCreds, dbCnfg)
		if err != nil {
			panic(err)
		}
		authUserServ, err := auth.NewAuthUser(*appCnfg, userRep)
		if err != nil {
			panic(err)
		}

		authUserRouter := authuserapi.AuthUserRouter{}
		authUserRouter.Init(apiGroup, authUserServ)
		// ---------------------
	}
	// Employee
	{
		// ----- Employee Auth -----
		employeeRep, err := employeerep.NewEmployeeRep(ctx, pgCreds, dbCnfg)
		if err != nil {
			panic(err)
		}
		authEmployeeServ, err := auth.NewAuthEmployee(*appCnfg, employeeRep)
		if err != nil {
			panic(err)
		}
		authEmployeeRouter := authemployeeapi.AuthEmployeeRouter{}
		authEmployeeRouter.Init(apiGroup, authEmployeeServ)
		// ---------------------

		// ----- Employee work -----
		employeeGroup := apiGroup.Group("/employee")
		employeeGroup.Use(middleware.AuthMiddleware(authEmployeeServ))

		artworkRep, err := artworkrep.NewArtworkRep(ctx, pgCreds, dbCnfg)
		if err != nil {
			panic(err)
		}
		artworkServ := artworkserv.NewArtworkService(artworkRep)
		artworkRouter := artworksapi.ArtworksRouter{}
		artworkRouter.Init(employeeGroup, artworkServ)
		// -------------------------
	}
	// ----- Admin Auth -----
	adminRep, err := adminrep.NewAdminRep(ctx, pgCreds, dbCnfg)
	if err != nil {
		panic(err)
	}
	authAdminServ, err := auth.NewAuthAdmin(*appCnfg, adminRep)
	if err != nil {
		panic(err)
	}
	authAdminRouter := authadminapi.AuthAdminRouter{}
	authAdminRouter.Init(apiGroup, authAdminServ)
	// ----------------------

	// ----- Admin work -----
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(authAdminServ))
	// ----------------------
	engine.Run(":8080")

	// server, err := api.NewServer()
	// if err != nil {
	// 	log.Fatal("cannot create server:", err)
	// }
	// err = server.Start(":8080")
	// if err != nil {
	// 	log.Fatal("cannot start server:", err)
	// }

	// r := gin.Default()
	// // route для Swagger - НЕ ТРОГАТЬ
	// url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// r.GET("/hello", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "Привет, мир! Вау"})
	// })

	// r.Run(":8080")
}

func main1() {
	wd, err := os.Getwd() // Получает директорию, из которой запущен `go run`
	if err != nil {
		panic(err)
	}
	fmt.Println("Working directory:", wd)

	ctx := context.Background()
	pgTestConfig := &cnfg.PostgresTestConfig{
		DbName:       "testArtwork",
		Port:         5432,
		Username:     "testUser",
		Password:     "testPassword",
		Image:        "postgres:latest",
		MigrationDir: "../../migrations/",
	}
	container, pgCreds, err := pgtest.NewTestPostgres(ctx, pgTestConfig)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer container.Terminate(ctx)
	fmt.Printf("Creds: %+v\n", pgCreds)

	if err = pgtest.MigrateUp(ctx, pgTestConfig.MigrationDir, &pgCreds); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	dbCnfg, err := cnfg.LoadDatebaseConfig("../../configs")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	urep, err := userrep.NewPgUserRep(ctx, &pgCreds, dbCnfg)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
	}
	users, err := urep.GetAll(ctx)
	if err != nil {
		if errors.Is(err, userrep.ErrUserNotFound) {
			fmt.Printf("No user\n")
		} else {
			fmt.Printf("ERROR: %v\n\n\n", err)
		}
		return
	}
	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}
	// newUser, _ := models.NewUser(
	// 	uuid.New(),
	// 	"test-user",
	// 	"test-login",
	// 	"hashed-password",
	// 	time.Now(),
	// 	"user@test.com",
	// 	true,
	// )
	// err = urep.Add(ctx, &newUser)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }

	// user, err := urep.GetByLogin(ctx, "test-login")
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// fmt.Printf("UPDATED %+v\n", *user)
}

func main2() {
	// redisCnfg, err := cnfg.LoadRedisCredentials()
	// if err != nil {
	// 	fmt.Printf("cannot load config: %v", err)
	// 	return
	// }
	// appCnfg, err := cnfg.LoadAppConfig()
	// if err != nil {
	// 	fmt.Printf("cannot load config: %v", err)
	// 	return
	// }
	// ctx := context.Background()
	// txRep, err := buyticketstxrep.NewBuyTicketsTxRep(ctx, redisCnfg)
	// if err != nil {
	// 	fmt.Printf("ERROR1: %v\n", err)
	// 	return
	// }
	// ticketServ, err := buyticketserv.NewBuyTicketsServ(txRep, *appCnfg)
	// if err != nil {
	// 	fmt.Printf("ERROR2: %v\n", err)
	// 	return
	// }
	// eventID := uuid.New()
	// cntTickets := 2
	// customerName := "customer1"
	// customerEmail := "customer@test.ru"
	// ttx, err := ticketServ.BuyTicket(ctx, eventID, cntTickets, customerName, customerEmail)
	// if err != nil {
	// 	fmt.Printf("ERROR3: %v\n", err)
	// 	return
	// }
	// err = ticketServ.ConfirmBuyTicket(ctx, ttx.GetID())
	// if err != nil {
	// 	fmt.Printf("ERROR4: %v\n", err)
	// 	return
	// }

	// // Config ------
	// pgTestCnfg, err := cnfg.LoadPgTestConfig()
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }
	// fmt.Printf("%+v\n", pgTestCnfg)
	// pgCreds, err := cnfg.LoadPgCredentials()
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }
	// dbCnfg, err := cnfg.LoadDatebaseConfig("./configs/")
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }
	// appCnfg, err := cnfg.LoadAppConfig()
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }
	// fmt.Printf("Postgres config: %+v\n", pgCreds)
	// fmt.Printf("App config: %+v\n", appCnfg)
	// fmt.Printf("Datebase config: %+v\n", dbCnfg)
	// // ------

	// // Repo ------
	// ctx := context.Background()
	// artrep, err := artworkrep.NewPgArtworkRep(ctx, pgCreds, dbCnfg)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// }
	// a, err := artrep.GetByID(ctx, uuid.MustParse("30154661-36c5-4761-96ea-691abb9bb407"))
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// }
	// eventrep, err := eventrep.NewPgEventRep(ctx, pgCreds, dbCnfg)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// }

	// startDate := time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC)
	// endDate := time.Date(2025, 6, 22, 0, 0, 0, 0, time.UTC)

	// events, err := eventrep.GetEventsOfArtworkOnDate(ctx, a, startDate, endDate)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// for _, a := range events {
	// 	fmt.Printf("%+v\n", a)
	// }

	// newUser, _ := models.NewUser(
	// 	uuid.New(),
	// 	"test-user",
	// 	"test-login",
	// 	"hashed-password",
	// 	time.Now(),
	// 	"user@test.com",
	// 	true,
	// )
	// user, err := urep.GetByLogin(ctx, "test-login")
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// err = urep.Delete(ctx, user.GetID())
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// err = urep.Add(ctx, user)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// updateFunc := func(u *models.User) (*models.User, error) {
	// 	newUser, err := models.NewUser(
	// 		u.GetID(),
	// 		"NEW USERNAME 10",
	// 		u.GetLogin(),
	// 		u.GetHashedPassword(),
	// 		u.GetCreatedAt(),
	// 		u.GetEmail(),
	// 		u.IsSubscribedToMail(),
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &newUser, nil
	// }
	// _, err = urep.Update(ctx, user.GetID(), updateFunc)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// err = urep.UpdateSubscribeToMailing(ctx, user.GetID(), false)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// user, err = urep.GetByLogin(ctx, "test-login")
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// 	return
	// }
	// fmt.Printf("UPDATED %+v\n", *user)

	// ------

	// // Repo ------
	// ures, err := userrep.NewPgUserRep(ctx, config)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n\n\n", err)
	// }
	// users, err := ures.GetAll(ctx)
	// if err != nil {
	// 	if errors.Is(err, userrep.ErrNoUser) {
	// 		fmt.Printf("No user\n")
	// 	} else {
	// 		fmt.Printf("ERROR: %v\n\n\n", err)
	// 	}
	// 	return
	// }
	// for _, user := range users {
	// 	fmt.Printf("%+v\n", user)
	// }
	// // ------

}
