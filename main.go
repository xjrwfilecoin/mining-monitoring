package main

import (
	"flag"
	"log"
	"mining-monitoring/app"
)

var configPath string
var workerPath string

func main() {
	flag.StringVar(&configPath, "configPath", "./configtest.json", "please config file configPath ")
	flag.StringVar(&workerPath, "workerPath", "./workerhost.json", "please config file workerPath ")
	flag.Parse()
	if configPath == "" {
		configPath = "./configtest.json"
	}
	if workerPath == "" {
		workerPath = "./workerhost.json"
	}
	err := app.Run(configPath)
	if err != nil {
		log.Fatal(err)
	}
}
