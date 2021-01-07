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
	err := app.Run(configPath,workerPath)
	if err != nil {
		log.Fatal(err)
	}
}
