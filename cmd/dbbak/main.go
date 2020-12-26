package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os/exec"
)

// 备份数据库
func main() {
	timer := cron.New(cron.WithSeconds())
	daySpec := "0 0 02 * * ?"
	entryID, err := timer.AddFunc(daySpec, func() {
		bakDb()
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		timer.Stop()
		timer.Remove(entryID)
	}()
	timer.Start()
	select {}
}

func bakDb() {
	fmt.Println("start execute script ...")
	cmd := exec.Command("/bin/bash", "-c", "sh bakdb.sh")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("end execute script ...")
}
