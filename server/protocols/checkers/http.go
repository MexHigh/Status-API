package checkers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"status-api/protocols"
	"status-api/structs"
)

var errTooManyRedirects = errors.New("Too many redirects")

// HTTP -
type HTTP struct{}

func (HTTP) ValidateConfig(config *structs.ServiceConfig) error {
	return nil
}

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

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	if skipSSLVerify, ok := c.ProtocolConfig["skip_ssl_verify"].(bool); ok && skipSSLVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// prepare client and request
	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
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
	if err != nil {
		if tempErr, ok := err.(*url.Error); ok && tempErr.Timeout() { // if error is timeout
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "timeout",
			}, nil
		} else if errors.Is(err, errTooManyRedirects) { // if error is errTooManyRedirects
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "too many redirects",
			}, nil
		} else if strings.Contains(err.Error(), "no route to host") { // TODO there might be a better way to solve this
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "no route to host",
			}, nil
		}
		return structs.CheckResult{}, err // unknown error
	}
	defer resp.Body.Close()

	if successCodes != nil {
		for _, sc := range successCodes {
			if resp.StatusCode == int(sc.(float64)) {
				return structs.CheckResult{
					Status: structs.Up,
					URL:    c.FriendlyURL,
				}, nil
			}
		}
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return structs.CheckResult{
				Status: structs.Up,
				URL:    c.FriendlyURL,
			}, nil
		}
	}

	return structs.CheckResult{
		Status: structs.Down,
		URL:    c.FriendlyURL,
		Reason: fmt.Sprintf("status code %d does not match conditions", resp.StatusCode),
	}, nil

}

func init() {
	protocols.Register("http", HTTP{})
}
