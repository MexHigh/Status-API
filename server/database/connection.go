package database

// A Database contains information about
// a database connection
type Database struct {
	Connected bool
}

// NewSQLite3 creates a new connection struct
// to an SQLite3 Database. Call db.Connect()
// to astablish the connection
func NewSQLite3(path string) *Database {
	return &Database{}
}

// Connect establishes the Database connection
func (db *Database) Connect() {

}

// Disconnect closes the connection
func (db *Database) Disconnect() {

}

// Exec executes a query on the database
func (db *Database) Exec(query string) interface{} {
	return nil
}
