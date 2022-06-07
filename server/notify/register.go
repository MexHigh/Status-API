package notify

import (
	"encoding/json"
	"fmt"
	"log"
	"status-api/structs"
	"time"
)

// Notifier defines an interface that can send messages to
// report a service as up or down.
//
// Each function will only be executed if the status is new.
// E.g. NotifyDown will only be called if the service was up
// already and vice versa.
type Notifier interface {
	NotifyDown(serviceName string, reportedDownAt time.Time, reason string) error
	NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration) error
}

// ConfigurableNotifier defines an extention for Notifier.
//
// It allows the Status API to provide the Notifier with a
// configuration from the config file at "notifiers > [name] > ."
//
// If this interface is implemented by a Notifier, it will be
// type asserted to ConfigurableNotifier automatically after
// the config file has been read. The Notifier is then provided
// with the raw JSON config, that he can unmarshal itself.
type ConfigurableNotifier interface {
	Notifier
	UnmarshalConfig(raw json.RawMessage) error
}

var (
	// notifiers holds all registed notifiers
	notifiers = make(map[string]Notifier)

	// activeNotifiers contains pointers to notifiers
	// explicitly activated through the config file
	activeNotifiers = make(map[string]*Notifier)
)

// Register registers a Notifier (or optionally
// a ConfigurableNotifier).
func Register(name string, notifier Notifier) {
	notifiers[name] = notifier
}

// GetNotifier returns the registered notifier, or nil,
// if it does not exist or has not been registered yet.
//
// This function also returns inactive notifiers.
func GetNotifier(notifier string) Notifier {
	if c, ok := notifiers[notifier]; ok {
		return c
	}
	return nil
}

// GetAllNotifierNames returns a list of notifier names
// after they have been registered. Otherwise the list
// will be empty.
//
// Set activeOnly to true to only list checkers, that
// were explicitly activated through the config file.
func GetAllNotifierNames(activeOnly bool) (names []string) {
	names = make([]string, 0, len(notifiers))
	if activeOnly {
		for key := range activeNotifiers {
			names = append(names, key)
		}
	} else {
		for key := range notifiers {
			names = append(names, key)
		}
	}
	return
}

// ProvideConfig provides the corresponding config for every
// Notifier that implements the ConfigurableNotifier interface.
//
// Notifiers that do not implement the ConfigurableNotifier will
// be skipped silently. Sucessfull config provisoning will be logged.
//
// The following JSON path will be provided as raw JSON:
// "notifiers > [name] > ."
func ProvideConfig(c *structs.Config) error {
	for key, notifier := range activeNotifiers {
		if cNotifier, ok := (*notifier).(ConfigurableNotifier); ok {
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

// Activate takes any number of notifier names to activate them.
// Returns an non-nil error if the function encounters a notifier
// name, that has not been registered.
func Activate(names ...string) error {
	for _, name := range names {
		if notifier, ok := notifiers[name]; ok {
			activeNotifiers[name] = &notifier
		} else {
			return fmt.Errorf("couldn't activate unregistered notifier '%s'", name)
		}
	}
	return nil
}
