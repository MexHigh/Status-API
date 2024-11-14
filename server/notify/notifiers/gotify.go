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
	defaultPriorty = 9
)

type gotifyConfig struct {
	Host                string `json:"host"`
	ApplicationKey      string `json:"application_key"`
	PriorityRaw         *int   `json:"priority,omitempty"`
	Priority            int    `json:"-"`
	CustomPriorityTimes []struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Priority int    `json:"priority"`
	} `json:"custom_priority_times"`
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
	// check if within custom priority timeslot
	prio, err := g.customPriority()
	if err != nil {
		return err
	}
	if prio == -1 {
		// not within custom priority timeslot
		prio = g.config.Priority // use "global" priority
	}

	// compose request body
	reqBodyMap := map[string]interface{}{
		"title":    title,
		"message":  message,
		"priority": prio,
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

// customPriority returns the priority if the current time is
// within a configured custom time (custom_priority_times prop).
// If not, it returns -1 and nil.
//
// If an error is returned, ommit the integer.
func (g *Gotify) customPriority() (int, error) {
	now := time.Now()
	for _, timeObj := range g.config.CustomPriorityTimes {
		from, err := time.ParseInLocation("15:04", timeObj.From, time.Local)
		if err != nil {
			return -1, err
		}
		from = from.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)

		to, err := time.ParseInLocation("15:04", timeObj.To, time.Local)
		if err != nil {
			return -1, err
		}
		to = to.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)

		if now.After(from) && now.Before(to) {
			return timeObj.Priority, nil
		}
	}
	return -1, nil
}

func (g *Gotify) NotifyDown(serviceName string, reportedDownAt time.Time, reason string) error {
	title := fmt.Sprintf(downNotificationTitle, serviceName)
	msg := fmt.Sprintf(downNotificationMsg, reportedDownAt.Local().Format(dateTimeFormat), reason)
	if err := g.gotifySend(title, msg); err != nil {
		return err
	}
	return nil
}

func (g *Gotify) NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration) error {
	title := fmt.Sprintf(upNotificationTitle, serviceName)
	msg := fmt.Sprintf(upNotificationMsg, reportedDownAt.Local().Format(dateTimeFormat), wasDownFor.String())
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
	if c.PriorityRaw != nil {
		c.Priority = *c.PriorityRaw
	} else {
		c.Priority = defaultPriorty
	}
	g.config = c
	return nil
}

// Interface guard (CAUTION: This interface guard does
// not detect, if a required function is implemented for
// the reciever type, which is invalid!)
var (
	_ notify.Notifier             = (*Gotify)(nil)
	_ notify.ConfigurableNotifier = (*Gotify)(nil)
)

// Register notifier
func init() {
	notify.Register("gotify", &Gotify{})
}
