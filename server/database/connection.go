package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Con is an instance of a
// database connection
var Con *gorm.DB

// InitializeSQLite3 creates a new connection struct
// to an SQLite3 Database. Call db.Connect()
// to astablish the connection
func InitializeSQLite3(dbPath string, dst ...interface{}) error {
	// open
	db, err := gorm.Open(
		sqlite.Open(dbPath),
		&gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		return err
	}

	log.Println("Connected to database at", dbPath)
	log.Println("Running database migrations")

	// auto migration
	if err := db.AutoMigrate(dst...); err != nil {
		return err
	}

	Con = db
	return nil

}
