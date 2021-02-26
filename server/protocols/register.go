package protocols

import (
	"status-api/schedules"
	"status-api/structs"
)

// Register registers a new protocol checker
// It can be used by adding it's name to the "protocol"
// key in the config
func Register(name string, checker structs.Checker) {
	schedules.Checkers[name] = checker
}
