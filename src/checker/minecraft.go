package checker

import (
	"net"
	"strconv"
	"strings"
	"time"

	"status-api/config"
	"status-api/minepong"
)

func checkMinecraft(name string, endpoint config.EndpointConfig) error {
	conn, err := net.DialTimeout(
		"tcp", 
		endpoint.MinecraftConfig.URL,
		time.Duration(5*time.Second),
	)
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			Status[name] = map[string]string{
				"url":    endpoint.FriedlyURL,
				"status": "down",
			}
		} else { // for example "no such host"
			return err
		}
		return nil
	}
	pong, err := minepong.Ping(conn, endpoint.MinecraftConfig.URL)
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
