package checker

import (
	"strings"
	"strconv"

	"status-api/config"
	"status-api/checker/minepong"
)

func checkMinecraft(name string, endpoint config.EndpointConfig) error {

	protocolConfig := endpoint.Protocol.Config.(*config.MinecraftConfig)

	pong, err := minepong.Ping(protocolConfig.URL)
	if err != nil {
		if e := err.Error(); strings.Contains(e, "i/o timeout") || strings.Contains(e, "connection refused") {
			Status[name] = map[string]string{
				"url":    endpoint.FriedlyURL,
				"status": "down",
			}
			return nil
		}
		return err // unknown error
	}

	Status[name] = map[string]string{
		"url":     endpoint.FriedlyURL,
		"status":  "up",
		"players": strconv.Itoa(pong.Players.Online) + "/" + strconv.Itoa(pong.Players.Max),
	}
	return nil
}
