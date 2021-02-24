package main

import (
	"flag"
	"fmt"
	"log"

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

	fmt.Println(
		c.ForService("Nextcloud").Protocol,
	)

}
