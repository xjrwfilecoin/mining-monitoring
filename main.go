package main

import (
	"flag"
	"log"
	"mining-monitoring/app"
)

var configPath string


func main() {
	flag.StringVar(&configPath, "configPath", "./config.json", "please config file configPath ")
	flag.Parse()
	err := app.Run(configPath)
	if err != nil {
		log.Fatal(err)
	}
}
