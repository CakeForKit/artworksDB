package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/pgtest"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
)

func main() {
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

	if err = pgtest.MigrateUp(ctx); err != nil {
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

func main1() {
	// Config ------
	pgTestCnfg, err := cnfg.LoadPgTestConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	fmt.Printf("%+v\n", pgTestCnfg)
	pgCreds, err := cnfg.LoadPgCredentials()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dbCnfg, err := cnfg.LoadDatebaseConfig("./configs/")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	appCnfg, err := cnfg.LoadAppConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	fmt.Printf("Postgres config: %+v\n", pgCreds)
	fmt.Printf("App config: %+v\n", appCnfg)
	fmt.Printf("Datebase config: %+v\n", dbCnfg)
	// ------

	// Repo ------
	ctx := context.Background()
	urep, err := userrep.NewPgUserRep(ctx, pgCreds, dbCnfg)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
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
	user, err := urep.GetByLogin(ctx, "test-login")
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
		return
	}
	err = urep.Delete(ctx, user.GetID())
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
		return
	}
	err = urep.Add(ctx, user)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
		return
	}
	updateFunc := func(u *models.User) (*models.User, error) {
		newUser, err := models.NewUser(
			u.GetID(),
			"NEW USERNAME 10",
			u.GetLogin(),
			u.GetHashedPassword(),
			u.GetCreatedAt(),
			u.GetEmail(),
			u.IsSubscribedToMail(),
		)
		if err != nil {
			return nil, err
		}
		return &newUser, nil
	}
	_, err = urep.Update(ctx, user.GetID(), updateFunc)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
		return
	}
	err = urep.UpdateSubscribeToMailing(ctx, user.GetID(), false)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
		return
	}
	user, err = urep.GetByLogin(ctx, "test-login")
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
		return
	}
	fmt.Printf("UPDATED %+v\n", *user)
	// err = urep.Add(ctx, newUser)
	// if err != nil {
	// 	fmt.Printf("Expected add ERROR: %v\n\n\n", err)
	// 	return
	// }

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
