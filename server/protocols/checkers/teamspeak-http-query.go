package checkers

import (
	"errors"
	"fmt"
	"strings"

	"status-api/protocols"
	"status-api/structs"

	ts "github.com/jkoenig134/go-ts3"
)

type TeamspeakHTTPQuery struct{}

func (TeamspeakHTTPQuery) ValidateConfig(c *structs.ServiceConfig) error {
	// get config params (returns err if something mandatory is required)
	baseURL, apiKey, err := getParams(c)
	if err != nil {
		return err
	}
	// check if the params works
	tsClient := ts.NewClient(ts.NewConfig(baseURL, apiKey))
	if _, err = tsClient.ServerInfo(); err != nil {
		return errors.New(getTSHTTPQueryErrorReason(err))
	}
	// everything setup correctly
	return nil
}

func (TeamspeakHTTPQuery) Check(name string, c *structs.ServiceConfig) (structs.CheckResult, error) {
	res := structs.CheckResult{
		URL: c.FriendlyURL,
	}

	// get config params
	baseURL, apiKey, err := getParams(c)
	if err != nil {
		return res, err
	}

	// define teamspeak query client
	tsClient := ts.NewClient(ts.NewConfig(baseURL, apiKey))

	// get server status (returns current users, max users, etc.)
	status, err := tsClient.ServerInfo()
	if err != nil {
		res.Status = structs.Down
		res.Reason = getTSHTTPQueryErrorReason(err)
		return res, nil
	}

	// at this point, we can be sure that the server is online
	res.Status = structs.Up

	// add additional infos
	res.Misc = map[string]string{
		"Clients Online": fmt.Sprintf("%d/%d", status.VirtualserverClientsonline-status.VirtualserverQueryclientsonline, status.VirtualserverMaxclients),
		"Latency":        fmt.Sprintf("%.1fms", status.VirtualserverTotalPing),
		"Uptime":         secondsToFiendlyTime(status.VirtualserverUptime),
	}

	return res, nil
}

// Register checker
func init() {
	protocols.Register("teamspeak-http", TeamspeakHTTPQuery{})
}

// getTSHTTPQueryErrorReason creates a meaningful error
// message from a TS HTTP Query Client error as a string.
func getTSHTTPQueryErrorReason(err error) string {
	e := err.Error()
	codeString := func(i int) string {
		return fmt.Sprintf(`"code":%d`, i)
	}

	// match the error codes
	if strings.Contains(e, codeString(5120)) {
		return "Query API key out of scope"
	}
	// TODO add more?

	// match all other errors
	if strings.Contains(e, "Error Code: 502") {
		return "Query unreachable"
	}

	return err.Error() // return the error if not matched
}

// getParams gets all protocol_config params (baseURL and apiKey).
// Returns an error, if anything that is required was not provided.
func getParams(c *structs.ServiceConfig) (string, string, error) {
	var baseURL string
	var apiKey string

	if bu, ok := c.ProtocolConfig["base_url"].(string); ok {
		baseURL = bu
	} else {
		// no base_url was provided, try friendly url
		baseURL = strings.Replace(c.FriendlyURL, "ts3server://", "https://", 1)
	}

	if ak, ok := c.ProtocolConfig["api_key"].(string); ok {
		apiKey = ak
	} else {
		return "", "", fmt.Errorf("protocol config field api_key is required")
	}

	return baseURL, apiKey, nil
}

// secondsToFiendlyTime converts seconds into `[DDd] HHh MMmin`
func secondsToFiendlyTime(seconds int) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %02dh %02dmin", days, hours, minutes)
	}
	return fmt.Sprintf("%02dh %02dmin", hours, minutes)
}
