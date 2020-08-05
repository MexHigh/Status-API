package main

import (
	"log"
	"flag"

	"status-api/api"
	"status-api/checker"
	"status-api/config"
)

func main() {

	configPath := flag.String("config", "./config.json", "Path to the config.json file")
	flag.Parse()

	log.Println("Loading endpoints from", *configPath)
	var err error
	config.Endpoints, err = config.LoadEndpointsFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Retrieving initial service states")
	checker.CheckAllServices()
	go checker.Updater(600) // 10 Minutes

	api.Start()

}
