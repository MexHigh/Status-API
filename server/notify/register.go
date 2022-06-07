package notify

import (
	"encoding/json"
	"log"
	"status-api/structs"
	"time"
)

// Notifier defines an interface that can send messages to
// report a service as up or down
type Notifier interface {
	NotifyDown(serviceName string, reportedDownAt time.Time, reason string)
	NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration)
}

type ConfigurableNotifier interface {
	Notifier
	UnmarshalConfig(raw json.RawMessage) error
}

var notifiers = make(map[string]Notifier)

func Register(name string, notifier Notifier) {
	notifiers[name] = notifier
}

// GetNotifier returns the registered checkers, or nil,
// if it does not exist or has not been registered yet
func GetNotifier(notifier string) Notifier {
	if c, ok := notifiers[notifier]; ok {
		return c
	}
	return nil
}

func GetAllNotifierNames() (names []string) {
	names = make([]string, 0, len(notifiers))
	for key := range notifiers {
		names = append(names, key)
	}
	return
}

// ProvideConfig provides the corresponding config for every
// Notifier that implements the ConfigurableNotifier interface.
//
// The following JSON path will be provided: notifiers > [name] > .
func ProvideConfig(c *structs.Config) error {
	for key, notifier := range notifiers {
		if cNotifier, ok := notifier.(ConfigurableNotifier); ok {
			config, ok := c.Notifiers[key]
			if !ok {
				log.Printf("Notifier '%s' did not recieve a configuration (not found in config file)", key)
				continue
			}
			if err := cNotifier.UnmarshalConfig(config); err != nil {
				return err
			}
			log.Printf("Notifier '%s' recieved a configuration", key)
		}
	}
	return nil
}
