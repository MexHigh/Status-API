package checkers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"status-api/structs"
	"time"
)

// HTTP -
type HTTP struct{}

var errTooManyRedirects = errors.New("Too many redirects")

// Check -
func (HTTP) Check(name string, c *structs.ServiceConfig) (structs.Result, error) {

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

	// do check
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
		return structs.Result{}, err
	}
	// TODO add credentials
	resp, err := client.Do(req)
	if tempErr, ok := err.(*url.Error); ok && tempErr.Timeout() { // if error is timeout
		return structs.Result{
			Status: "down",
			URL:    c.FriendlyURL,
			Reason: "timeout",
		}, nil
	} else if errors.Is(err, errTooManyRedirects) {
		return structs.Result{
			Status: "down",
			URL:    c.FriendlyURL,
			Reason: "too many redirects",
		}, nil
	} else if err != nil {
		return structs.Result{}, err // TODO is this ok to return an empty result?
	}

	if successCodes != nil {
		for _, sc := range successCodes {
			if resp.StatusCode == int(sc.(float64)) {
				return structs.Result{
					Status: "up",
					URL:    c.FriendlyURL,
				}, nil
			}
		}
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return structs.Result{
				Status: "up",
				URL:    c.FriendlyURL,
			}, nil
		}
	}

	return structs.Result{
		Status: "down",
		URL:    c.FriendlyURL,
		Reason: fmt.Sprintf("status code %d does not match conditions", resp.StatusCode),
	}, nil

}

// Interface guard
var _ structs.Checker = (*HTTP)(nil)
