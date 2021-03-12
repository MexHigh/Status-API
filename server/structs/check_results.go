package structs

import (
	"time"
)

// Status is a string enumaration of a
// service status ("up"|"down")
type Status string

const ( // TODO this is currently unused !!!
	// Up = "up"
	Up Status = "up"
	// Down = "down"
	Down Status = "down"
	// Problems = "problems"
	Problems Status = "problems"
)

// A CheckResult is the status, url and an optional
// downtime reason at a specific point in time
type CheckResult struct {
	Status Status `json:"status"`
	URL    string `json:"url"`
	Reason string `json:"reason,omitempty"`
}

// CheckResults wraps multiple CheckResults structs
// by service name combined with a timestamp
type CheckResults struct {
	At       time.Time              `json:"at"`
	Services map[string]CheckResult `json:"services"`
}
