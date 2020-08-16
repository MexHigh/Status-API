package checker

import (
	"strconv"
	"net/http"
	
	"status-api/config"
)

func checkHTTP(name string, endpoint config.EndpointConfig) error {
	r, err := http.Get(endpoint.URL)
	if err != nil {
		return err
	}
	for _, statusCode := range endpoint.SuccessOn {
		if r.StatusCode == statusCode {
			Status[name] = map[string]string{
				"url":    endpoint.URL,
				"status": "up",
				"code":   strconv.Itoa(r.StatusCode),
			}
			return nil
		}
	}
	Status[name] = map[string]string{
		"url":    endpoint.URL,
		"status": "down",
		"code":   strconv.Itoa(r.StatusCode),
	}
	return nil
}