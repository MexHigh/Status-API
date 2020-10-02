package checker

import (
	"strconv"
	"strings"
	"status-api/checker/minepong"
)

// MinecraftConfig is the config struct for testing if a minecraft server is up.
// It also displays the number of online players in the API response.
type MinecraftConfig struct {
	URL string `json:"url"`
}

func (m *MinecraftConfig) setDefaults() {
	return
}

func checkMinecraft(name string, endpoint *Endpoint) error {

	protocolConfig := endpoint.Protocol.Config.(*MinecraftConfig)

	pong, err := minepong.Ping(protocolConfig.URL)
	if err != nil {
		if e := err.Error(); strings.Contains(e, "i/o timeout") || strings.Contains(e, "connection refused") {
			endpoint.Status = map[string]string{
				"url":    endpoint.FriedlyURL,
				"status": "down",
			}
			return nil
		}
		return err // unknown error
	}

	endpoint.Status = map[string]string{
		"url":     endpoint.FriedlyURL,
		"status":  "up",
		"players": strconv.Itoa(pong.Players.Online) + "/" + strconv.Itoa(pong.Players.Max),
	}
	return nil
}