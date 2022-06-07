package notify

import (
	"status-api/structs"
	"time"
)

var servicesThatAreDown = make(map[string]*structs.CheckResultWithNameAndTime, 0)

func ReportDown(result *structs.CheckResultWithNameAndTime) {
	if service, ok := servicesThatAreDown[result.Name]; ok {
		// means that the service is already down --> do nothing
		return
	} else {
		// means that the service was up before --> notify
		// ad add to servicesThatAreDown
		notifyDown(service, time.Now())
		servicesThatAreDown[service.Name] = service
	}
}

func notifyDown(result *structs.CheckResultWithNameAndTime, now time.Time) {
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
		notifyUp(service, time.Now())
		delete(servicesThatAreDown, service.Name)
	} // else: means that service is already up --> do nothing
}

func notifyUp(result *structs.CheckResultWithNameAndTime, now time.Time) {
	for _, notifier := range notifiers {
		go func(n Notifier, r *structs.CheckResultWithNameAndTime) {
			n.NotifyUp(r.Name, r.Time, time.Since(r.Time))
		}(notifier, result)
	}
}
