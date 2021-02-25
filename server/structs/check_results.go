package structs

import (
	"time"
)

// Status is a string enumaration of a
// service status ("up"|"down")
type Status string

const (
	// Up = "up"
	Up Status = "up"
	// Down = "down"
	Down Status = "down"
	// Problems = "problems"
	Problems Status = "problems"
)

type CheckResult struct {
	Status Status `json:"status"`
	URL    string `json:"url"`
	Reason string `json:"reason,omitempty"`
}

type CheckResults struct {
	At       time.Time              `json:"at"`
	Services map[string]CheckResult `json:"services"`
}
