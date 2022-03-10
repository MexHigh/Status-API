package checkers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"status-api/protocols"
	"status-api/structs"
)

var errTooManyRedirects = errors.New("too many redirects")

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
	} else { // use friendy url if "test_url" is not specified
		testURL = c.FriendlyURL
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	// skip sslVerify check
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

	// do the request
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

	// check success code
	if successCodes, ok := c.ProtocolConfig["success_codes"].([]interface{}); ok && len(successCodes) > 0 {
		sawMatch := false
		for _, sc := range successCodes {
			if resp.StatusCode == int(sc.(float64)) {
				sawMatch = true
			}
		}
		if !sawMatch {
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: fmt.Sprintf("status code %d did not match conditions", resp.StatusCode),
			}, nil
		}
	} else { // just check if the status code is between 200-299 if not specified
		if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: fmt.Sprintf("status code %d did not match conditions", resp.StatusCode),
			}, nil
		}
	}

	// check expected content
	if expContent, ok := c.ProtocolConfig["exptected_content"].(string); ok {
		// parse HTTP response to string
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return structs.CheckResult{}, err
		}
		respString := string(respBytes)

		if !strings.Contains(respString, expContent) {
			return structs.CheckResult{
				Status: structs.Down,
				URL:    c.FriendlyURL,
				Reason: "response did not contain expected string",
			}, nil
		} // else continue
	}

	return structs.CheckResult{
		Status: structs.Up,
		URL:    c.FriendlyURL,
	}, nil

}

func init() {
	protocols.Register("http", HTTP{})
}
