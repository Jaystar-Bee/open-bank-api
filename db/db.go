package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var MainDB *sql.DB

func InitDatabase() {
	Database, err := sql.Open("sqlite3", "api.db")
	if err != nil {
		panic(err)
	}
	MainDB = Database
	MainDB.SetMaxOpenConns(14)
	MainDB.SetMaxIdleConns(7)
	CreateTables()
}
