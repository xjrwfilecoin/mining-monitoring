package main

import (
	"fmt"
	"os/exec"
)

func main() {
	for i:=0;i<100;i++{
		cmd := exec.Command("sshpass", "-p", "", "ssh", "-p", "22", "root@worker", "free", "-h")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("error: ", err.Error())
			return
		}
		fmt.Println(string(output))
	}
}
