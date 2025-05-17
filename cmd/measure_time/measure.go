package main

import (
	"flag"
	"fmt"

	querytime "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/query_time"
)

func main() {
	start, stop, step := 1000, 51000, 2000

	flag.Parse()
	args := flag.Args()
	if !(len(args) > 0 && args[0] == "g") {
		measurer, err := querytime.NewQueryTime()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		err = measurer.MeasureTime(start, stop, step, true)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	err := querytime.DrawGraph(start, stop, step)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
