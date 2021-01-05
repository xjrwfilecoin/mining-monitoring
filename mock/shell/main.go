package main

import (
	"fmt"
	"mining-monitoring/shellParsing"
)

func main() {
	manager, err := shellParsing.NewManager("./workerhost.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result, err := manager.DoShell()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("result: %v  \n", result)
}


