package structs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// This file defines an own model to be used
// with gorm

type Model struct {
	ID uint `gorm:"primary_key" json:"-"`
}

type CheckResultsModel struct {
	Model
	Data CheckResults
}

// Value implements the driver.Valuer interface for ResultMap
func (a CheckResults) Value() (driver.Value, error) {
	bytes, err := json.Marshal(&a)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Scan implements the sql.Scanner interface for ResultMap
func (a *CheckResults) Scan(value interface{}) error {
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
var _ driver.Valuer = (*CheckResults)(nil)
var _ sql.Scanner = (*CheckResults)(nil)

type ArchiveResultsModel struct {
	Model
	Data ArchiveResults
}

// Value implements the driver.Valuer interface for ResultMap
func (a ArchiveResults) Value() (driver.Value, error) {
	bytes, err := json.Marshal(&a)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Scan implements the sql.Scanner interface for ResultMap
func (a *ArchiveResults) Scan(value interface{}) error {
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
var _ driver.Valuer = (*ArchiveResults)(nil)
var _ sql.Scanner = (*ArchiveResults)(nil)
