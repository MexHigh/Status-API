package notifiers

import (
	"log"
	"time"

	"status-api/notify"
)

// Log is a simple notifier that just logs
// notifications to stdout via log.Println
type Log struct{}

func (Log) NotifyDown(serviceName string, reportedDownAt time.Time, reason string) error {
	log.Printf("Service '%s' is DOWN! Reason: %s", serviceName, reason)
	return nil
}

func (Log) NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration) error {
	log.Printf("Service '%s' is UP again! Was down since %s for %s", serviceName, reportedDownAt.Local().Format(dateTimeFormat), wasDownFor.String())
	return nil
}

// Interface guard
var _ notify.Notifier = (*Log)(nil)

// Register notifier
func init() {
	notify.Register("log", Log{})
}
