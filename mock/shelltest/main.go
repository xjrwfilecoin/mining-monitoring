package main

import (
	"fmt"
	"os/exec"
)

func main() {



}


// 网络不通
func test02() {
	// error:  exit status 255
	cmd := exec.Command("sshpass", "-p", "", "ssh", "-p", "22521", "root@ya-node1000", "free", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	fmt.Println(string(output))
}

// 免密失效
func test01() {
	// 免密失效 error:  exit status 6
	cmd := exec.Command("sshpass", "-p", "", "ssh", "-p", "22521", "root@ya-node100", "free", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	fmt.Println(string(output))
}
