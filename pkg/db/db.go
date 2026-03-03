package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var err error

func Init(dbFile string) error {
	var install bool

	_, err := os.Stat(dbFile)

	if err != nil {
		install = true
	}

	schema := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '', 
		title CHAR(16) NOT NULL DEFAULT '',
		comment CHAR(256) NOT NULL DEFAULT '',
		repeat CHAR(128) NOT NULL DEFAULT ''
	);
	CREATE INDEX index_date ON scheduler (date);
	`

	db, err = sql.Open("sqlite", dbFile)

	if err != nil {
		return err
	}
	//defer db.Close()
	if install {

		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			return err
		}
	}

	return nil
}
