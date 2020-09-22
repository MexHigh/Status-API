package checker

import (
	"strings"
	"strconv"
	"net/http"
	
	"status-api/config"
)

func checkHTTP(name string, endpoint config.EndpointConfig) error {
	// use friendly URL if test URL in HTTP Config is not set
	var testURL string
	if t := endpoint.HTTPConfig.TestURL; t == "" {
		testURL = endpoint.FriedlyURL
	} else {
		testURL = t
	}

	r, err := http.Get(testURL)
	if err != nil {
		return err
	}
	successCodeStrings := strings.Split(
		strings.ReplaceAll(endpoint.HTTPConfig.SuccessCodes, " ", ""),
		",",
	)
	for _, statusCodeString := range successCodeStrings {
		statusCode, err := strconv.Atoi(statusCodeString)
		if err != nil {
			return err
		}
		if r.StatusCode == statusCode {
			Status[name] = map[string]string{
				"url":    endpoint.FriedlyURL,
				"status": "up",
				"code":   strconv.Itoa(r.StatusCode),
			}
			return nil
		}
	}
	Status[name] = map[string]string{
		"url":    endpoint.FriedlyURL,
		"status": "down",
		"code":   strconv.Itoa(r.StatusCode),
	}
	return nil
}