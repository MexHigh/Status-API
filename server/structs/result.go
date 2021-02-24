package structs

import "time"

// Status is a string enumaration of a
// service status ("up"|"down")
type Status string

const (
	// Up = "up"
	Up Status = "up"
	// Down = "down"
	Down Status = "down"
)

// Misc is a map, that will be flattened into
// the Result struct when being marshalled
type Misc map[string]string // other informations

// A Result represents the health and some other
// information of a service after it was checked
type Result struct {
	Status Status `json:"status"`
	URL    string `json:"url"`
	Misc   `json:",omitempty"`
}

// Results wraps multiple Results by name and
// adds a timestamp at when they were checked
type Results struct {
	At       time.Time
	Services map[string]Result
}
