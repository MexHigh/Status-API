package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"status-api/database"
	"status-api/notify"
	_ "status-api/notify/notifiers" // enforce compilation of all notifiers
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

	checkers := strings.Join(protocols.GetAllCheckerNames(), ", ")
	log.Println("Loaded protocol checkers:", checkers)

	notifiers := strings.Join(notify.GetAllNotifierNames(false), ", ")
	log.Println("Loaded notifiers:", notifiers)

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

	log.Println("Activating notifiers mentioned in config")
	keys := make([]string, 0)
	for key := range c.Notifiers {
		keys = append(keys, key)
	}
	if err := notify.Activate(keys...); err != nil {
		panic(err)
	}

	activeNotifiers := strings.Join(notify.GetAllNotifierNames(true), ", ")
	log.Println("Activated notifiers:", activeNotifiers)

	log.Println("Providing config to activated notifiers")
	if err := notify.ProvideConfig(c); err != nil {
		panic(err)
	}

	log.Println("Connecting to SQLite3 database")
	if err := database.InitializeSQLite3(
		c.DBPath,
		&structs.CheckResultsModel{},
		&structs.ArchiveResultsModel{},
		&structs.AtomFeedItemModel{},
	); err != nil {
		panic(err)
	}

	log.Println("Starting trigger jobs")
	go schedules.StartCheckTriggerJob(c)
	go schedules.StartArchiveTriggerJob(c)

	log.Println("Starting server")
	if err := server.Start(
		c.Host,
		c.FrontendPath,
		c.DashboardTitle,
		c.DashboardLogoPath,
		!c.NoFrontend,
		c.AllowedAPIKeys,
	); err != nil {
		panic(err)
	}

}
