package checkers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"status-api/protocols"
	"status-api/structs"

	"github.com/Tnze/go-mc/bot"
)

// Minecraft -
type Minecraft struct{}

// Check -
func (Minecraft) Check(name string, c *structs.ServiceConfig) (structs.CheckResult, error) {

	var hostPort string
	if hp, ok := c.ProtocolConfig["server_address"].(string); ok {
		hostPort = hp
	} else {
		hostPort = c.FriendlyURL
		if !strings.Contains(hostPort, ":") {
			hostPort = hostPort + ":25565"
		}
	}

	var res = structs.CheckResult{
		URL: c.FriendlyURL,
	}

	resp, _, err := bot.PingAndListTimeout(hostPort, 10*time.Second)
	if err != nil {
		res.Status = structs.Down
		if e := err.Error(); strings.Contains(e, "i/o timeout") {
			res.Reason = "I/O timeout"
		} else if strings.Contains(e, "connection refused") {
			res.Reason = "connection refused"
		} else if strings.Contains(e, "no route to host") {
			res.Reason = "no route to host"
		} else if strings.Contains(e, "no such host") {
			res.Reason = "no such host"
		} else {
			res.Reason = e
		}
		return res, nil
	}

	var s status
	if err = json.Unmarshal(resp, &s); err != nil {
		return structs.CheckResult{}, err
	}

	res.Status = structs.Up
	res.Misc = map[string]string{
		"version":        s.Version.Name,
		"players_online": fmt.Sprintf("%d/%d", s.Players.Online, s.Players.Max),
	}

	return res, nil
}

// Register checker
func init() {
	protocols.Register("minecraft", Minecraft{})
}

type status struct {
	Players struct {
		Max    int
		Online int
	}
	Version struct {
		Name string
	}
}
