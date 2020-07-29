package main

import (
	"flag"
	"./api"
	"./config"
	"./checker"
)

func main() {

	configPath := flag.String("config", "../config.json", "Path to the config.json file")
	flag.Parse()

	var err error
	config.Conf, err = config.LoadConfigFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	checker.CheckAllServices()
	go checker.Updater(600) // 10 Minutes

	api.Start()

}