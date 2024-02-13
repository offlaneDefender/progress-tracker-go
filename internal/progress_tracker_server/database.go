package progress_tracker

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const file string = "progtrac.db"

const create string = `
	CREATE TABLE IF NOT EXISTS goals(
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT not null,
		progress INTEGER NOT NULL
	);
`

const insert string = `
	INSERT INTO goals(name, progress) VALUES("eepyyy", 0);
`

type ProgtracDB struct {
	db *sql.DB
}

func ConnectToDBAndInit() (*ProgtracDB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &ProgtracDB{
		db: db,
	}, nil

}

func InsertTestData(db *ProgtracDB) (int, error) {
	err := db.db.Ping()

	if err != nil {
		return 0, err
	}

	res, err := db.db.Exec(insert)

	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return int(id), nil
}
