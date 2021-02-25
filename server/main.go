package main

import (
	"flag"
	"fmt"
	"log"

	"status-api/database"
	"status-api/schedules"
	"status-api/server"
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
		&structs.CheckResultsModel{},
		&structs.ArchiveResultsModel{},
		// TODO add archiveStruct
	); err != nil {
		panic(err)
	}

	log.Println("Starting trigger routines")
	go schedules.CheckTriggerRoutine(c)
	go schedules.ArchiveTriggerRoutine(c)

	log.Println("Starting server")
	if err := server.Start(
		c.Host,
		!c.NoFrontend,
		c.FrontendPath,
	); err != nil {
		panic(err)
	}

}
