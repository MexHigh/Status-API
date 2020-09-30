package checker

import (
	"errors"
	"time"
	"net/url"
	"net/http"
)

var errTooManyRedirects = errors.New("Too many redirects")

// HTTPConfig is the config struct for testing if a website or API is reachable.
// It supports basicauth, too (via the embedded Credentials *struct).
type HTTPConfig struct {
	// If the StatusCodes were not set it defaults to 200
	SuccessCodes []int `json:"success_codes,omitempty"`
	// If the test URL is empty, the friendly URL will be used
	TestURL string `json:"test_url,omitempty"`
	// If the Credentials are set, basicauth will be used to
	// authenticate against the webserver. This is nil by default
	// as it is a struct pointer.
	Credentials *struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials,omitempty"`
}

func (h *HTTPConfig) setDefaults() {
	if len(h.SuccessCodes) == 0 {
		h.SuccessCodes = []int{200}
	}
}

func (h *HTTPConfig) Check(name string, endpoint Endpoint) error {

	protocolConfig := endpoint.Protocol.Config.(*HTTPConfig)

	// inline function to access endpoint var and set service status
	mark := func(status string) {
		endpoint.Status[name] = map[string]string{
			"url":    endpoint.FriedlyURL,
			"status": status,
		}
	}

	// use friendly URL if test URL in HTTP Config is not set
	var testURL string
	if t := protocolConfig.TestURL; t == "" {
		testURL = endpoint.FriedlyURL
	} else {
		testURL = t
	}

	// do the request
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
		return err
	}
	if creds := protocolConfig.Credentials; creds != nil {
		req.SetBasicAuth(creds.Username, creds.Password)
	}
	resp, err := client.Do(req)
	if tempErr, ok := err.(*url.Error); ok && tempErr.Timeout() { // if error is a timeout
		mark("down")
		return nil
	} else if errors.Is(err, errTooManyRedirects) { // if error is caused by too many redirects
		mark("down (too many redirects)")
		return nil
	} else if err != nil { // unknown error
		return err
	}

	// check for matching status code
	/*successCodeStrings := strings.Split(
		strings.ReplaceAll(endpoint.HTTPConfig.SuccessCodes, " ", ""),
		",",
	)*/
	for _, statusCode := range protocolConfig.SuccessCodes {
		if resp.StatusCode == statusCode {
			mark("up")
			return nil
		} else if (resp.StatusCode == 401) { // only matches if 401 was not set as "success_codes"
			mark("down (unauthorized)")
			return nil
		}
	}

	// if everything before this point has not returned
	// mark the service as down
	mark("down")
	return nil
}