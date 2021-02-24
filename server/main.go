package main

import (
	"flag"
	"log"

	"status-api/schedules"
	"status-api/structs"
)

var configPath = flag.String("config", "./config.json", "Path to the config.json file")

func main() {

	flag.Parse()

	log.Println("Loading endpoints from", *configPath)
	c, err := structs.ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Starting trigger routines")
	go schedules.CheckTriggerRoutine(c, 10)
	go schedules.ArchiveTriggerRoutine(c)

	log.Println("Starting API server")
	// TODO add API server
	for {
	}

}
