package checkers

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"status-api/protocols"
	"status-api/structs"
)

// SMTP -
type SMTP struct{}

// Check -
func (SMTP) Check(name string, c *structs.ServiceConfig) (structs.CheckResult, error) {

	// get query address from config
	var hostPort string
	if hp, ok := c.ProtocolConfig["server_address"].(string); ok {
		hostPort = hp
	} else {
		// TODO check for port if required by Dial
		hostPort = c.FriendlyURL
	}

	smtpClient, err := smtp.Dial(hostPort)
	if err != nil {
		return structs.CheckResult{
			Status: structs.Down,
			URL:    c.FriendlyURL,
			Reason: err.Error(),
		}, nil
	}
	defer smtpClient.Close()

	// DEBUG
	fmt.Printf("Starting noop")

	if err = smtpClient.Noop(); err != nil {
		return structs.CheckResult{
			Status: structs.Down,
			URL:    c.FriendlyURL,
			Reason: "NOOP failed: " + err.Error(),
		}, nil
	}

	// DEBUG
	fmt.Printf("noop end")

	// check if PlainAuth works if specified
	if plainAuth, ok := c.ProtocolConfig["plain_auth"].(map[string]interface{}); ok {
		var identity, username, password, host string

		if tempIdentity, ok := plainAuth["identity"].(string); ok {
			identity = tempIdentity
		} // else identity is empty string, which is expected
		if tempUsername, ok := plainAuth["username"].(string); ok {
			username = tempUsername
		} else {
			return structs.CheckResult{}, errors.New("smtp: username is required in \"plain_auth\"")
		}
		if tempPassword, ok := plainAuth["password"].(string); ok {
			password = tempPassword
		} else {
			return structs.CheckResult{}, errors.New("smtp: password is required in \"plain_auth\"")
		}
		if tempHost, ok := plainAuth["host"].(string); ok {
			host = tempHost
		} else {
			hostOnly := strings.Split(hostPort, ":")[0]
			host = hostOnly
		}

		// DEBUG
		fmt.Printf("Using:\n%s\n%s\n%s\n%s\n",
			identity, username, password, host,
		)

		auth := smtp.PlainAuth(identity, username, password, host)
		if err := smtpClient.Auth(auth); err != nil {
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "SMTP auth error: " + err.Error(),
			}, nil
		}
	}

	// report up
	return structs.CheckResult{
		Status: structs.Up,
		URL:    c.FriendlyURL,
	}, nil

}

// Register checker
func init() {
	protocols.Register("smtp", SMTP{})
}
