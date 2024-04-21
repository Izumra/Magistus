package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(
	driverName string,
	dbPath string,
) *Storage {
	db, err := sql.Open(driverName, dbPath)
	if err != nil {
		panic("Occured the error while connecting to the db " + err.Error())
	}

	return &Storage{
		db,
	}
}
