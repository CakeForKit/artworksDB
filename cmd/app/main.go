package main

import (
	"context"
	"fmt"
	"log"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/config"
)

func main() {
	// Config ------
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	fmt.Printf("Postgres config: %+v\n", config.Postgres)
	fmt.Printf("App config: %+v\n", config.App)
	fmt.Printf("Datebase config: %+v\n", config.Datebase)
	// ------

	// Repo ------
	ctx := context.Background()
	urep, err := userrep.NewPgUserRep(ctx, config)
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
			u.GetMail(),
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
