package checker

import (
	"time"
	"net/url"
	"strings"
	"strconv"
	"net/http"
	
	"status-api/config"
)

func checkHTTP(name string, endpoint config.EndpointConfig) error {

	// inline function to access endpoint and set service status
	mark := func(status string) {
		Status[name] = map[string]string{
			"url":    endpoint.FriedlyURL,
			"status": status,
		}
	}

	// use friendly URL if test URL in HTTP Config is not set
	var testURL string
	if t := endpoint.HTTPConfig.TestURL; t == "" {
		testURL = endpoint.FriedlyURL
	} else {
		testURL = t
	}

	// do the request
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return err
	}
	if creds := endpoint.HTTPConfig.Credentials; creds != nil {
		req.SetBasicAuth(creds.Username, creds.Password)
	}
	resp, err := client.Do(req)
	if err, ok := err.(*url.Error); ok && err.Timeout() { // if error is a timeout
		mark("down")
		return nil
	} else if err != nil { // other error
		return err
	}

	// check for matching status code
	successCodeStrings := strings.Split(
		strings.ReplaceAll(endpoint.HTTPConfig.SuccessCodes, " ", ""),
		",",
	)
	for _, statusCodeString := range successCodeStrings {
		statusCode, err := strconv.Atoi(statusCodeString)
		if err != nil {
			return err
		}
		if resp.StatusCode == statusCode {
			mark("up")
			return nil
		} else if (resp.StatusCode == 401) {
			mark("unauthorized")
			return nil
		}
	}

	// if everything before this point has not returned
	// mark the service as down
	mark("down")
	return nil
}