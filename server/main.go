package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"status-api/database"
	"status-api/protocols"
	_ "status-api/protocols/checkers" // enforce compilation of all checkers
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

	names := "" + strings.Join(protocols.GetAllCheckerNames(), ", ")
	log.Println("Loaded protocol checkers:", names)

	log.Println("Loading config from", *configPath)
	c, err := structs.ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Validating config")
	if err := protocols.ValidateConfig(c); err != nil {
		panic(err)
	} else {
		log.Println("Config is valid, continuing")
	}

	log.Println("Connecting to SQLite3 database")
	if err := database.InitializeSQLite3(
		c.DBPath,
		&structs.CheckResultsModel{},
		&structs.ArchiveResultsModel{},
	); err != nil {
		panic(err)
	}

	log.Println("Starting trigger jobs")
	go schedules.StartCheckTriggerJob(c)
	go schedules.StartArchiveTriggerJob(c)

	log.Println("Starting server")
	if err := server.Start(
		c.Host,
		!c.NoFrontend,
		c.FrontendPath,
	); err != nil {
		panic(err)
	}

}
