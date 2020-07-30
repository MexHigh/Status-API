package main

import (
	"flag"

	"./api"
	"./checker"
	"./config"
)

func main() {

	configPath := flag.String("config", "../config.json", "Path to the config.json file")
	flag.Parse()

	var err error
	config.Endpoints, err = config.LoadEndpointsFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	checker.CheckAllServices()
	go checker.Updater(600) // 10 Minutes

	api.Start()

}
