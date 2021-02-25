package structs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
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

// A Result represents the health and some other
// information of a service after it was checked
type Result struct {
	Status Status `json:"status"`
	URL    string `json:"url"`
	Reason string `json:"reason,omitempty"`
}

// ResultMap maps a result to it's service name
type ResultMap map[string]Result

// Value implements the driver.Valuer interface for ResultMap
func (r ResultMap) Value() (driver.Value, error) {
	bytes, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Scan implements the sql.Scanner interface for ResultMap
func (r *ResultMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Error casting to byte slice")
	}
	if err := json.Unmarshal(bytes, r); err != nil {
		return err
	}
	return nil
}

// Interface guards for ResultMap
var _ driver.Valuer = (ResultMap)(nil)
var _ sql.Scanner = (*ResultMap)(nil)

// Justification: ResultMap must implement the Valuer and Scanner
// interfaces as a map[string]Result cannot be saved to the database.
// Those functions are used to marshal the map[string]Result to JSON
// (Scan to retreive from Database, Value to store to Database)

// Results wraps multiple Results by name and
// adds a timestamp at when they were checked
type Results struct {
	Model
	At       time.Time `json:"at"`
	Services ResultMap `json:"services"`
}
