package notifiers

import (
	"encoding/json"
	"fmt"
	"time"

	"status-api/notify"
)

const (
	downNotificationTitle = "Service '%s' is down!"
	downNotificationMsg   = "Reported at: %s.\nReason: %s"
	upNotificationTitle   = "Service '%s' is up again!"
	upNotificationMsg     = "Reported at: %s (was down for %s)"
)

type Gotify struct{}

func (Gotify) NotifyDown(serviceName string, reportedDownAt time.Time, reason string) {

}

func (Gotify) NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration) {

}

type gotifyConfig struct {
	Key string `json:"key"`
}

func (Gotify) UnmarshalConfig(raw json.RawMessage) error {
	var c gotifyConfig
	if err := json.Unmarshal(raw, &c); err != nil {
		return err
	}
	fmt.Println(c.Key)
	return nil
}

// Register notifier
func init() {
	notify.Register("gotify", Gotify{})
}
