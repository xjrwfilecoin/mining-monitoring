package main

import (
	"net/http"
	_ "net/http/pprof"
)
func main() {

	err := http.ListenAndServe(":7070", nil)
	if err!=nil{
		panic(err)
	}
}
