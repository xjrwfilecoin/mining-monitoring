package main

import (
	"encoding/json"
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
	bytes, err := json.Marshal(result)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(bytes))
}


