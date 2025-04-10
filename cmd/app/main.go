package main

import (
	"fmt"
	"log"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/util"
)

func main() {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	fmt.Printf("%+v", config)
	fmt.Println("It works!")
}
