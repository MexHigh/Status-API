package checker

import (
	ts3 "github.com/multiplay/go-ts3"
)

// TSConfig is the config struct for testing if a Teamspeak 3/5 server is online
type TSConfig struct {
	QueryURL string `json:"query_url"`
}

func (t *TSConfig) setDefaults() {
	return
}

func checkTeamspeak(name string, endpoint *Endpoint) error {

	protocolConfig := endpoint.Protocol.Config.(*TSConfig)

	tsclient, err := ts3.NewClient(protocolConfig.QueryURL)
	if err != nil {
		// on failure
		endpoint.Status = map[string]string{
			"url":    endpoint.FriedlyURL,
			"status": "down",
		}
		return nil
	}
	defer tsclient.Close()
	
	// on success
	endpoint.Status = map[string]string{
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