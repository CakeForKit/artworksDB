package main

import (
	"fmt"

	querytime "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/query_time"
)

func main() {
	measurer, err := querytime.NewQueryTime()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	err = measurer.MeasureTime()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
