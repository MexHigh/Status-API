package checker

import (
	"status-api/config"

	ts3 "github.com/multiplay/go-ts3"
)

func checkTeamspeak(name string, endpoint config.EndpointConfig) error {

	protocolConfig := endpoint.Protocol.Config.(*config.TSConfig)

	tsclient, err := ts3.NewClient(protocolConfig.QueryURL)
	if err != nil {
		// on failure
		Status[name] = map[string]string{
			"url":    endpoint.FriedlyURL,
			"status": "down",
		}
		return nil
	}
	defer tsclient.Close()
	
	// on success
	Status[name] = map[string]string{
		"url":    endpoint.FriedlyURL,
		"status": "up",
	}
	return nil

	/*if err := tsclient.Login(endpoint.Credentials.Username, endpoint.Credentials.Password); err != nil {
		panic(err)
	}

	if v, err := tsclient.Version(); err != nil {
		panic(err)
	} else {
		log.Println("Server is running version:", v)
	}*/
	
}