package structs

// This file defines an own model to be used
// with gorm

type Model struct {
	ID uint `gorm:"primary_key" json:"-"`
}
