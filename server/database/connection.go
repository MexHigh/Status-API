package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Con is an instance of a
// database connection
var Con *gorm.DB

// InitializeSQLite3 creates a new connection struct
// to an SQLite3 Database. Call db.Connect()
// to astablish the connection
func InitializeSQLite3(path string, dst ...interface{}) error {

	// open
	db, err := gorm.Open(
		sqlite.Open(path),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}

	// auto migration
	if err := db.AutoMigrate(dst...); err != nil {
		return err
	}

	Con = db
	return nil

}
