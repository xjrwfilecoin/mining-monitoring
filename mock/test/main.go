package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {

	for i := 0; i < 500; i++ {
		go func() {
			cmd := exec.Command("sshpass", "-p", "", "ssh", "root@node01", "free", "-h")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(output))
			time.Sleep(5 * time.Second)

		}()
	}

	time.Sleep(1 * time.Hour)
}
