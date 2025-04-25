package main

import (
	"context"
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/userrep"
)

func main() {
	// config, err := util.LoadConfig("../../")
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }
	// var id uuid.UUID = uuid.New()
	// id = uuid.Nil
	// fmt.Printf("%v", id)
	// fmt.Printf("%+v", config)
	// fmt.Println("It works!")

	ctx := context.Background()
	ures, err := userrep.NewPgUserRep(ctx)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
	}
	users, err := ures.GetAll(ctx)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n\n", err)
	}
	for _, user := range users {
		fmt.Printf("%+v", user)
	}
	// fmt.Printf("%v\n\n\n", ures.TestSelect(ctx))

	// ctx := context.Background()
	// connCtx, cancelConn := context.WithTimeout(ctx, 5*time.Second)
	// defer cancelConn()
	// queryCtx, cancelQuery := context.WithTimeout(ctx, 3*time.Second)
	// defer cancelQuery()

	// connStr := "postgres://puser:ppassword@postgres_container:5432/artworks"
	// config, err := pgxpool.ParseConfig(connStr)
	// if err != nil {
	// 	log.Fatal("Ошибка конфигурации пула: ", err)
	// }
	// dbpool, err := pgxpool.NewWithConfig(connCtx, config)
	// if err != nil {
	// 	log.Fatal("Ошибка создания пула: ", err)
	// }
	// defer dbpool.Close()

	// sql := "select username, email from users"
	// rows, sqlerr := dbpool.Query(queryCtx, sql)
	// if sqlerr != nil {
	// 	panic(fmt.Sprintf("QueryRow failed: %v", sqlerr))
	// }

	// for rows.Next() {
	// 	var username string
	// 	var email string
	// 	rows.Scan(&username, &email)
	// 	fmt.Printf("%s\t%s\n\n\n", username, email)
	// 	log.Printf("%s\t%s\n\n\n", username, email)
	// 	// log.In("Ошибка создания пула: ", err)
	// }
}
