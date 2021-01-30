package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("sshpass", "-p", "", "ssh", "root@worker01", "free", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error: ",err.Error())
		return
	}
	fmt.Println(string(output))
}
