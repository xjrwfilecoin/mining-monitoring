package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("this is a error：%s  %s \n", r, time.Now())
		}
	}()
	cmd := exec.Command("/bin/bash", "-c", "sh deploy.sh")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("exec shell %v \n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("hello"))
	fmt.Printf("received deploy msg ... %v \n", time.Now())
	return

}

func main() {
	http.HandleFunc("/deploy", DeployHandler)
	err := http.ListenAndServe(":8899", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
