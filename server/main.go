package main

import (
	"flag"
	"fmt"
	"log"

	"status-api/api"
	"status-api/database"
	"status-api/schedules"
	"status-api/structs"
)

var configPath = flag.String("config", "./config.json", "Path to the config.json file")

func main() {

	fmt.Println()
	fmt.Println("  ~ Status API by leon.wtf ~  ")
	fmt.Println()

	flag.Parse()

	log.Println("Loading config from", *configPath)
	c, err := structs.ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Connecting to SQLite3 database")
	if err := database.InitializeSQLite3(
		c.DBPath,
		&structs.Results{},
		// TODO add archiveStruct
	); err != nil {
		panic(err)
	}

	log.Println("Starting trigger routines")
	go schedules.CheckTriggerRoutine(c, c.CheckInterval)
	go schedules.ArchiveTriggerRoutine(c)

	log.Println("Starting API server")
	if err := api.Start(c.APIHost); err != nil {
		panic(err)
	}

}
