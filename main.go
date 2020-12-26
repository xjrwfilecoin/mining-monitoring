package main

import (
	"log"
	"mining-monitoring/app"
)

func main() {
	err := app.Run("./configtest.json")
	if err != nil {
		log.Fatal(err)
	}
}
