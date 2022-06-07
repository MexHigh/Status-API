package notify

import (
	"log"
	"time"

	"status-api/structs"
)

var servicesThatAreDown = make(map[string]*structs.CheckResultWithNameAndTime, 0)

// ReportDown can be called every time a service was tested unreachable.
// It checks itself if the service was reported down previously to not
// send multiple notifications
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
	for notifierName, notifier := range activeNotifiers {
		go func(nName string, n Notifier, r *structs.CheckResultWithNameAndTime) {
			err := n.NotifyDown(r.Name, r.Time, r.Result.Reason)
			if err != nil {
				log.Printf("Error in notifier '%s': %s", nName, err.Error())
			}
		}(notifierName, *notifier, result)
	}
}

// ReportUp can be called every time a service was tested reachable.
// It checks itself if the service was reported up previously to not
// send multiple notifications
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
	for notifierName, notifier := range activeNotifiers {
		go func(nName string, n Notifier, r *structs.CheckResultWithNameAndTime) {
			err := n.NotifyUp(r.Name, r.Time, time.Since(r.Time))
			if err != nil {
				log.Printf("Error in notifier '%s': %s", nName, err.Error())
			}
		}(notifierName, *notifier, result)
	}
}
