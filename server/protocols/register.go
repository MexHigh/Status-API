package protocols

import (
	"fmt"
	"log"
	"status-api/structs"
)

// Checker defines a struct that can perform protocol-specific checks
type Checker interface {
	Check(name string, config *structs.ServiceConfig) (structs.CheckResult, error)
}

// ValidatableChecker allows you to define a ValidateConfig function for
// your checker. This function is invoked right after the config has been
// loaded and before the checking functions are scheduled. This is optional
// but advised as wrong configurations will be detected before invoking
// any connections.
type ValidatableChecker interface {
	ValidateConfig(config *structs.ServiceConfig) error
}

// checkers will be filled by protocols.Register()
// called from within a checker to register itself
var checkers = make(map[string]Checker)

// Register registers a new protocol checker.
// A registered checker can be used by adding it's
// name to the "protocol" key in the config.
// The checker struct must implement the Checker
// interface and return a CheckResult and an error.
// If errors are not handled inside the Check() method and
// are returned, the service will be flagged as down and the
// error description is passed to the response in the "reason" key.
func Register(name string, checker Checker) {
	checkers[name] = checker
}

// GetChecker returns the registered checkers, or nil,
// if it does not exist or has not been registered yet
func GetChecker(protocol string) Checker {
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

// ValidateConfig validates the parsed configuration
// against the registered config checkers that implement the
// ValidatableConfig interface
func ValidateConfig(config *structs.Config) error {
	for name, config := range config.Services {
		// check if the checker exists for this service config
		if c := GetChecker(config.Protocol); c != nil {
			// check if the corresponding checker implements ValidatableChecker
			if vc, ok := c.(ValidatableChecker); ok {
				// perform the validation
				if err := vc.ValidateConfig(&config); err != nil {
					return fmt.Errorf("could not validate service config for service \"%s\": %s", name, err.Error())
				} // else nothing, validation was successfull
			} else {
				log.Printf("Skipping config validation for service \"%s\" (checker \"%s\" cannot validate the configuration)", name, config.Protocol)
			}
		} else {
			return fmt.Errorf("protocol %s not supported", config.Protocol)
		}
	}
	return nil
}
