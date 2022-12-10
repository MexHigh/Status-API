package checkers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

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
		hostPort = c.FriendlyURL
		// append default port if not specified
		hostPortSplit := strings.Split(hostPort, ":")
		if len(hostPortSplit) == 1 || hostPortSplit[1] == "" {
			hostPort += ":587"
		}
	}

	// create the connection and smtp client
	conn, err := net.DialTimeout("tcp", hostPort, 5*time.Second)
	if err != nil {
		return structs.CheckResult{
			Status: structs.Down,
			URL:    c.FriendlyURL,
			Reason: err.Error(),
		}, nil
	}
	defer conn.Close()
	smtpClient, err := smtp.NewClient(conn, c.FriendlyURL)
	if err != nil {
		fmt.Println("Penis")
		return structs.CheckResult{
			Status: structs.Down,
			URL:    c.FriendlyURL,
			Reason: err.Error(),
		}, nil
	}
	defer smtpClient.Close()

	// execute a NOOP operation to check if the server is responding
	if err = smtpClient.Noop(); err != nil {
		return structs.CheckResult{
			Status: structs.Down,
			URL:    c.FriendlyURL,
			Reason: "NOOP failed: " + err.Error(),
		}, nil
	}

	// upgrade with StartTLS if wanted
	if useStartTLS, ok := c.ProtocolConfig["use_starttls"].(bool); ok && useStartTLS {
		// check if STARTTLS extension is supported
		if supported, _ := smtpClient.Extension("STARTTLS"); !supported {
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "STARTTLS extension not supported",
			}, nil
		}

		// upgrade
		if err := smtpClient.StartTLS(&tls.Config{
			ServerName: c.FriendlyURL,
		}); err != nil {
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "StartTLS error: " + err.Error(),
			}, nil
		}
	}

	// check if PlainAuth works if specified
	if plainAuth, ok := c.ProtocolConfig["plain_auth"].(map[string]interface{}); ok {
		// check if AUTH extension is supported
		if supported, _ := smtpClient.Extension("AUTH"); !supported {
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "AUTH extension not supported",
			}, nil
		}

		// provision value from config
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

		// do the authentication
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
