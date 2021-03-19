package checkers

import (
	"strings"

	"golang.org/x/crypto/ssh"

	"status-api/protocols"
	"status-api/structs"
)

// TeamspeakSSHQuery -
type TeamspeakSSHQuery struct{}

// Check -
func (TeamspeakSSHQuery) Check(name string, c *structs.ServiceConfig) (structs.CheckResult, error) {

	// get query address from config
	var hostPort string
	if hp, ok := c.ProtocolConfig["query_address"].(string); ok {
		hostPort = hp
	} else {
		hostPort = strings.ReplaceAll(c.FriendlyURL, "ts3server://", "")
		hostPort = strings.Split(hostPort, "?")[0]
		hostPort = hostPort + ":10022"
	}

	var res = structs.CheckResult{
		URL: c.FriendlyURL,
	}

	sshConfig := &ssh.ClientConfig{
		/* User: "serveradmin"
		Auth: []ssh.AuthMethod{
			ssh.Password(""),
		},*/
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	_, err := ssh.Dial("tcp", hostPort, sshConfig)
	if e := err.Error(); strings.Contains(e, "connection refused") || strings.Contains(e, "no route to host") { // ssh endpoint not responding
		res.Status = structs.Down
	} else if strings.Contains(e, "unable to authenticate") { // authentication error, but reachable
		res.Status = structs.Up
	} else if err != nil { // unknown error
		res.Status = structs.Down
		res.Reason = e
	} else { // connection could be established???
		res.Status = structs.Down
		res.Reason = "Connection could be established without authentication! Please check your teamspeak config!"
	}

	return res, nil

}

// Register checker
func init() {
	protocols.Register("teamspeak-ssh", TeamspeakSSHQuery{})
}
