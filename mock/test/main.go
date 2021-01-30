package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("test")
	})
	err := http.ListenAndServe(":8899", nil)
	if err!=nil{
		panic(err)
	}
}
