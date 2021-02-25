package main

import (
	"flag"
	"log"

	"status-api/database"
	"status-api/schedules"
	"status-api/structs"
)

var configPath = flag.String("config", "./config.json", "Path to the config.json file")
var dbPath = flag.String("database", "./db.sqlite", "Path to the SQLite3 database")

func main() {

	flag.Parse()

	log.Println("Loading config from", *configPath)
	c, err := structs.ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Connecting to SQLite3 database at", *dbPath)
	if err := database.InitializeSQLite3(
		*dbPath,
		&structs.Results{},
		// TODO add archiveStruct
	); err != nil {
		panic(err)
	}

	log.Println("Starting trigger routines")
	//go schedules.CheckTriggerRoutine(c, 10)
	go schedules.ArchiveTriggerRoutine(c)

	log.Println("Starting API server")
	// TODO add API server
	for {
	}

}
