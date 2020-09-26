package checker

import (
	"net"
	"strconv"
	"strings"
	"time"

	"status-api/config"
	"status-api/checker/minepong"
)

func checkMinecraft(name string, endpoint config.EndpointConfig) error {

	protocolConfig := endpoint.Protocol.Config.(*config.MinecraftConfig)

	conn, err := net.DialTimeout(
		"tcp", 
		protocolConfig.URL,
		time.Duration(5*time.Second),
	)
	if err != nil {
		if e := err.Error(); strings.Contains(e, "i/o timeout") || strings.Contains(e, "connection refused") {
			Status[name] = map[string]string{
				"url":    endpoint.FriedlyURL,
				"status": "down",
			}
		} else {
			return err
		}
		return nil
	}
	pong, err := minepong.Ping(conn, protocolConfig.URL)
	if err != nil {
		return err
	}
	Status[name] = map[string]string{
		"url":     endpoint.FriedlyURL,
		"status":  "up",
		"players": strconv.Itoa(pong.Players.Online) + "/" + strconv.Itoa(pong.Players.Max),
	}
	return nil
}
