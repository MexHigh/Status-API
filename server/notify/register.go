package notify

import "time"

// Notifier defines an interface that can send messages to
// report a service as up or down
type Notifier interface {
	NotifyDown(serviceName string, reportedDownAt time.Time, reason string)
	NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration)
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
