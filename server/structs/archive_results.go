package structs

import (
	"time"
)

// A Downtime combines a timestamp and a reason of
// when and why a service was not reachable
type Downtime struct {
	At     time.Time `json:"at"`
	Reason string    `json:"reason,omitempty"`
}

type ArchiveResult struct {
	Status       string     `json:"status"`
	Availability float64    `json:"availability"`
	Downtimes    []Downtime `json:"downtimes"`
}

type ArchiveResults struct {
	At       time.Time                `json:"at"`
	Services map[string]ArchiveResult `json:"services"`
}
