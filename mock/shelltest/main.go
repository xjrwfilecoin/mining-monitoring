package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("lotus-miner", "actor","control","list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error: ",err.Error())
		return
	}
	fmt.Println(string(output))
}
