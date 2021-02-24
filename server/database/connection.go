package database

type DBConnection struct {
	Connected bool
}

func (db *DBConnection) Connect() {

}

func (db *DBConnection) Disconnect() {

}

func (db *DBConnection) Exec(query string) interface{} {

}
