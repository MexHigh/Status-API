package structs

import (
	"time"
)

// Status is a string enumaration of a
// service status ("up"|"problems"|"down")
type Status string

const (
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
	Status Status            `json:"status"`
	URL    string            `json:"url"`
	Reason string            `json:"reason,omitempty"`
	Misc   map[string]string `json:"misc,omitempty"`
}

// CheckResults wraps multiple CheckResults structs
// by service name combined with a timestamp
type CheckResults struct {
	At       time.Time              `json:"at"`
	Services map[string]CheckResult `json:"services"`
}
