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
	result := manager.DoShell()
	fmt.Printf("result: %v  \n", result)
}

func MinerInfo() {
	resMap := make(map[string]interface{})
	shellParse := shellParsing.NewShellParse()

	err := shellParse.GetMinerInfo(resMap)
	if err != nil {
		fmt.Printf("minerInfo error %v \n", err)
		return
	}
	fmt.Printf("minerInfo:  %v \n", resMap)

	err = shellParse.GetPostBalance(resMap)
	if err != nil {
		fmt.Printf("postBalance error %v \n", err)
		return
	}
	fmt.Printf("postBalance:  %v \n", resMap)

	err = shellParse.GetMinerJobs(resMap)
	if err != nil {
		fmt.Printf("minerJobs error %v \n", err)
		return
	}

	fmt.Printf("minerJobs: %v \n", resMap)

	hardwareInfo, err := shellParse.BatchHardwareInfo()
	if err != nil {
		fmt.Printf("hardwareInfo:  %v \n", err)
		return
	}
	fmt.Printf("hardwareInfo: %v \n", hardwareInfo)

}
