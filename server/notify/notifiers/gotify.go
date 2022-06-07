package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"status-api/notify"
)

const (
	downNotificationTitle = "Service '%s' is down!"
	downNotificationMsg   = "Reported at: %s\nReason: %s"
	upNotificationTitle   = "Service '%s' is up again!"
	upNotificationMsg     = "Reported at: %s\nWas down for: %s)"

	defaultPriorty = 9
)

type gotifyConfig struct {
	Host                 string `json:"host"`
	ApplicationKey       string `json:"application_key"`
	OverwritePriorityRaw *int   `json:"overwrite_priority,omitempty"`
	Priority             int    `json:"-"`
}

// Gotify uses a Gotify server to send
// notifications via WebSocket
//
// See https://github.com/gotify/server
//
// Gotify implements notify.ConfigurableNotifier
type Gotify struct {
	config gotifyConfig
}

func (g *Gotify) gotifySend(title, message string) error {
	// compose request body
	reqBodyMap := map[string]interface{}{
		"title":    title,
		"message":  message,
		"priority": g.config.Priority,
	}
	reqBody, err := json.Marshal(reqBodyMap)
	if err != nil {
		return err
	}

	// client and request
	client := http.DefaultClient
	req, err := http.NewRequest("POST", g.config.Host+"/message", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Gotify-Key", g.config.ApplicationKey)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// check some response data
	if res.StatusCode != 200 {
		return fmt.Errorf("got status %d instead of 200", res.StatusCode)
	}
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var respBodyMap map[string]interface{}
	if err := json.Unmarshal(respBody, &respBodyMap); err != nil {
		return err
	}
	if _, ok := respBodyMap["id"]; !ok {
		return fmt.Errorf("expected field 'id' is not in response")
	}

	return nil
}

func (g *Gotify) NotifyDown(serviceName string, reportedDownAt time.Time, reason string) error {
	title := fmt.Sprintf(downNotificationTitle, serviceName)
	msg := fmt.Sprintf(downNotificationMsg, reportedDownAt.Local().String(), reason)
	if err := g.gotifySend(title, msg); err != nil {
		return err
	}
	return nil
}

func (g *Gotify) NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration) error {
	title := fmt.Sprintf(upNotificationTitle, serviceName)
	msg := fmt.Sprintf(upNotificationMsg, reportedDownAt.Local().String(), wasDownFor.String())
	if err := g.gotifySend(title, msg); err != nil {
		return err
	}
	return nil
}

func (g *Gotify) UnmarshalConfig(raw json.RawMessage) error {
	var c gotifyConfig
	if err := json.Unmarshal(raw, &c); err != nil {
		return err
	}
	if c.OverwritePriorityRaw != nil {
		c.Priority = *c.OverwritePriorityRaw
	} else {
		c.Priority = defaultPriorty
	}
	g.config = c
	return nil
}

// Interface guard (CAUTION: This interface guard does
// not detect, if a required function is implemented for
// the reciever type, which is invalid!)
var _ notify.ConfigurableNotifier = (*Gotify)(nil)

// Register notifier
func init() {
	notify.Register("gotify", &Gotify{})
}
