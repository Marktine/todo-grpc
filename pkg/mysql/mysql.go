package mysql

import (
	"database/sql"
)

// IDBManager interface
type IDBManager interface {
	Open (dataSourceName string) error
	Close ()
}

// Manager - mysql database manager
type Manager struct {
	DB *sql.DB 
}

// Open - open database connection
func (mdb *Manager) Open(dataSourceName string) error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	mdb.DB = db
	return nil
}

// Close - close the database connection
func (mdb *Manager) Close() {
	mdb.DB.Close()
}