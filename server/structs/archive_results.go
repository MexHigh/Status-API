package structs

import (
	"time"
)

// A Downtime combines a start and end timestamp
// and a reason of when and why a service was not reachable
type Downtime struct {
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	Reason string    `json:"reason,omitempty"`
}

// ArchiveResult represents the status and availability of
// a service over a longer period. Downtimes are tracked as well
type ArchiveResult struct {
	Status       Status     `json:"status"`
	Availability float64    `json:"availability"`
	Downtimes    []Downtime `json:"downtimes"`
}

// ArchiveResults wrapps multiple ArchiveResult structs
// by service name combined with a timestamp
type ArchiveResults struct {
	At       time.Time                `json:"at"`
	Services map[string]ArchiveResult `json:"services"`
}
