package structs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/gorilla/feeds"
)

// This file defines an own model to be used with gorm

// Model sets a primary key field to model structs
type Model struct {
	ID uint `gorm:"primary_key" json:"-"`
}

// CheckResultsModel wraps CheckResults and makes
// it JSON serializable in the database
type CheckResultsModel struct {
	Model
	Data CheckResults
}

// Value implements the driver.Valuer interface for CheckResultsModel
func (a CheckResults) Value() (driver.Value, error) {
	bytes, err := json.Marshal(&a)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Scan implements the sql.Scanner interface for CheckResultsModel
func (a *CheckResults) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error casting to byte slice")
	}
	if err := json.Unmarshal(bytes, a); err != nil {
		return err
	}
	return nil
}

// Interface guards for ResultMap
var _ driver.Valuer = (*CheckResults)(nil)
var _ sql.Scanner = (*CheckResults)(nil)

// ArchiveResultsModel wraps ArchiveResults and makes
// it JSON serializable in the database
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
		return errors.New("error casting to byte slice")
	}
	if err := json.Unmarshal(bytes, a); err != nil {
		return err
	}
	return nil
}

// Interface guards for ResultMap
var _ driver.Valuer = (*ArchiveResults)(nil)
var _ sql.Scanner = (*ArchiveResults)(nil)

// TODO
type AtomFeedItem feeds.Item

// TODO
type AtomFeedItemModel struct {
	Model
	Data AtomFeedItem
}

// Value implements the driver.Valuer interface for TODO
func (a AtomFeedItem) Value() (driver.Value, error) {
	bytes, err := json.Marshal(&a)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Scan implements the sql.Scanner interface for TODO
func (a *AtomFeedItem) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error casting to byte slice")
	}
	if err := json.Unmarshal(bytes, a); err != nil {
		return err
	}
	return nil
}

// Interface guards for ResultMap
var _ driver.Valuer = (*AtomFeedItem)(nil)
var _ sql.Scanner = (*AtomFeedItem)(nil)
