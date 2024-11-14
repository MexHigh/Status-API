package notifiers

const (
	dateTimeFormat = "02.01.2006, 15:04:05 -0700 MST"

	downNotificationTitle = "Service '%s' is down!"
	downNotificationMsg   = "Reported at: %s\nReason: %s"
	upNotificationTitle   = "Service '%s' is up again!"
	upNotificationMsg     = "Reported at: %s\nWas down for: %s)"
)
