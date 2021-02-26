package checkers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"status-api/protocols"
	"status-api/structs"
)

var errTooManyRedirects = errors.New("Too many redirects")

// HTTP -
type HTTP struct{}

// Check -
func (HTTP) Check(name string, c *structs.ServiceConfig) (structs.CheckResult, error) {

	// test url
	var testURL string
	if t, ok := c.ProtocolConfig["test_url"].(string); ok {
		testURL = t
	} else {
		testURL = c.FriendlyURL
	}

	// success codes
	var successCodes []interface{}
	if s, ok := c.ProtocolConfig["success_codes"].([]interface{}); ok {
		successCodes = s
	} // else remain nil

	// prepare client and request
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errTooManyRedirects
			}
			return nil
		},
	}
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return structs.CheckResult{}, err
	}

	// check for basicauth credentials
	if creds, ok := c.ProtocolConfig["credentials"].(map[string]interface{}); ok {
		req.SetBasicAuth(
			creds["username"].(string),
			creds["password"].(string),
		)
	}

	// do it
	resp, err := client.Do(req)
	if tempErr, ok := err.(*url.Error); ok && tempErr.Timeout() { // if error is timeout
		return structs.CheckResult{
			Status: "down",
			URL:    c.FriendlyURL,
			Reason: "timeout",
		}, nil
	} else if errors.Is(err, errTooManyRedirects) {
		return structs.CheckResult{
			Status: "down",
			URL:    c.FriendlyURL,
			Reason: "too many redirects",
		}, nil
	} else if err != nil {
		return structs.CheckResult{}, err // TODO is this ok to return an empty result?
	}

	if successCodes != nil {
		for _, sc := range successCodes {
			if resp.StatusCode == int(sc.(float64)) {
				return structs.CheckResult{
					Status: "up",
					URL:    c.FriendlyURL,
				}, nil
			}
		}
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return structs.CheckResult{
				Status: "up",
				URL:    c.FriendlyURL,
			}, nil
		}
	}

	return structs.CheckResult{
		Status: "down",
		URL:    c.FriendlyURL,
		Reason: fmt.Sprintf("status code %d does not match conditions", resp.StatusCode),
	}, nil

}

func init() {
	protocols.Register("http", HTTP{})
}
