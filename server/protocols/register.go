package protocols

import (
	"status-api/structs"
)

// checkers will be filled by protocols.Register()
// called from within a checker to register itself
var checkers = make(map[string]structs.Checker)

// Register registers a new protocol checker.
// A registered checker can be used by adding it's
// name to the "protocol" key in the config.
// The checker struct must implement the structs.Checker
// interface and return a CheckResult and an error.
// If errors are not handled inside the Check() method and
// are returned, the service will be flagged as down and the
// error description is passed to the response in the "reason" key.
func Register(name string, checker structs.Checker) {
	checkers[name] = checker
}

// GetChecker returns the registered checkers, or nil,
// if it does not exist or has not been registered yet
func GetChecker(protocol string) structs.Checker {
	if c, ok := checkers[protocol]; ok {
		return c
	}
	return nil
}

// GetAllCheckerNames returns a string slice containing
// all names of the checkers, that were already registered.
// They will be returned in the same format as it can be
// used in the config.json.
func GetAllCheckerNames() (names []string) {
	names = make([]string, 0, len(checkers))
	for key := range checkers {
		names = append(names, key)
	}
	return
}
