package main

import (
	"log"
	"flag"

	"status-api/api"
	"status-api/checker"
	"status-api/config"
)

var configPath = flag.String("config", "./config.json", "Path to the config.json file")

func main() {

	flag.Parse()

	log.Println("Loading endpoints from", *configPath)
	if err := config.LoadEndpointsFromFile(*configPath); err != nil {
		panic(err)
	}

	log.Println("Retrieving initial service states")
	go checker.Updater(600) // 10 Minutes

	api.Start()

}
