package notify

import (
	"time"

	"status-api/structs"
)

var servicesThatAreDown = make(map[string]*structs.CheckResultWithNameAndTime, 0)

func ReportDown(result *structs.CheckResultWithNameAndTime) {
	if _, ok := servicesThatAreDown[result.Name]; ok {
		// means that the service is already down --> do nothing
		return
	} else {
		// means that the service was up before --> notify
		// ad add to servicesThatAreDown
		notifyDown(result)
		servicesThatAreDown[result.Name] = result
	}
}

// notifyDown calls all notifier's NotifyDown function
func notifyDown(result *structs.CheckResultWithNameAndTime) {
	for _, notifier := range notifiers {
		go func(n Notifier, r *structs.CheckResultWithNameAndTime) {
			n.NotifyDown(r.Name, r.Time, r.Result.Reason)
		}(notifier, result)
	}
}

func ReportUp(result *structs.CheckResultWithNameAndTime) {
	if service, ok := servicesThatAreDown[result.Name]; ok {
		// means that the service was down before --> notify
		// and remove from servicesThatAreDown
		notifyUp(service)
		delete(servicesThatAreDown, service.Name)
	} // else: means that service is already up --> do nothing
}

// notifyUp calls all notifier's NotifyUp function
func notifyUp(result *structs.CheckResultWithNameAndTime) {
	for _, notifier := range notifiers {
		go func(n Notifier, r *structs.CheckResultWithNameAndTime) {
			n.NotifyUp(r.Name, r.Time, time.Since(r.Time))
		}(notifier, result)
	}
}
