package main

import (
	"fmt"
	"mining-monitoring/shellParsing"
	"time"
)

func main() {
	manager, err := shellParsing.NewManager()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sign := make(chan shellParsing.CmdData, 10)

	go func() {
		for {
			select {
			case info := <-sign:
				fmt.Println("result: ", info)
			}
		}
	}()
	manager.RunV1(sign)

	time.Sleep(60 * time.Hour)
}
