package notifiers

import (
	"log"
	"time"

	"status-api/notify"
)

// Log is a simple notifier that just logs
// notifications to stdout via log.Println
type Log struct{}

func (Log) NotifyDown(serviceName string, reportedDownAt time.Time, reason string) {
	log.Printf("Service '%s' is DOWN! Reason: %s", serviceName, reason)
}

func (Log) NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration) {
	log.Printf("Service '%s' is UP again! Was down since %s for %s", serviceName, reportedDownAt.Local().String(), wasDownFor.String())
}

// Register notifier
func init() {
	notify.Register("log", Log{})
}
