package main

import (
	"context"
	"fmt"
	"log"

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
	ures, err := userrep.NewPgUserRep(ctx, config)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
	}
	users, err := ures.GetAll(ctx)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
	}
	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}
	// ------

}
