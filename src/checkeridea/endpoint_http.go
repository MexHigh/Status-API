package checkeridea

import (
	"net/url"
	"errors"
	"time"
	"net/http"
)

type HTTPEndpoint struct {
	defaultEndpoint
	Config struct {
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
	} `json:"-"`
}

func (e *HTTPEndpoint) SetDefaults() {
	if e.Config.TestURL == "" {
		e.Config.TestURL = e.FriendlyURL
	}
	if len(e.Config.SuccessCodes) == 0 {
		e.Config.SuccessCodes = []int{200}
	}
}

var errTooManyRedirects = errors.New("Too many redirects")

func (e *HTTPEndpoint) Check() error {

	// inline function to access endpoint var and set service status
	mark := func(status string) {
		e.status = make(map[string]string)
		e.status = map[string]string{
			"url":    e.FriendlyURL,
			"status": status,
		}
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
	req, err := http.NewRequest("GET", e.Config.TestURL, nil)
	if err != nil {
		return err
	}
	if creds := e.Config.Credentials; creds != nil {
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
	for _, statusCode := range e.Config.SuccessCodes {
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

func (e *HTTPEndpoint) Status() EndpointStatus {
	return e.status
}