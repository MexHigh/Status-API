package structs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// A Downtime combines a timestamp and a reason of
// when and why a service was not reachable
type Downtime struct {
	At     time.Time `json:"at"`
	Reason string    `json:"reason,omitempty"`
}

// An ArchiveResult stores the cumulative state and availability
// of all Results for one day (or archiving period)
type ArchiveResult struct {
	Status       Status     `json:"status"`
	Availability float64    `json:"availability"`
	Downtimes    []Downtime `json:"downtimes"`
}

// An ArchiveResultMap maps an ArchiveResult to it's service name
type ArchiveResultMap map[string]ArchiveResult

// Value implements the driver.Valuer interface for ResultMap
func (a ArchiveResultMap) Value() (driver.Value, error) {
	bytes, err := json.Marshal(&a)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Scan implements the sql.Scanner interface for ResultMap
func (a *ArchiveResultMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Error casting to byte slice")
	}
	if err := json.Unmarshal(bytes, a); err != nil {
		return err
	}
	return nil
}

// Interface guards for ResultMap
var _ driver.Valuer = (ArchiveResultMap)(nil)
var _ sql.Scanner = (*ArchiveResultMap)(nil)

// ArchiveResults wraps an ArchiveResultsMap and a timestamp
// which describes the Date of the archiving
type ArchiveResults struct {
	Model
	At       time.Time        `json:"at"`
	Services ArchiveResultMap `json:"services"`
}
