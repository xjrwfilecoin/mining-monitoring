package main

import (
	"flag"
	"log"
	"mining-monitoring/app"
)

var path string

func main() {
	flag.StringVar(&path,"path","./configtest.json","please config file path ")
	flag.Parse()
	err := app.Run(path)
	if err != nil {
		log.Fatal(err)
	}
}
