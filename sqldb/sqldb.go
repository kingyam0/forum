package sqldb

import "database/sql"

// DB is a global variable to hold db connection
var DB *sql.DB

// ConnectDB opens a connection to the database
func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		panic(err.Error())
	}

	DB = db

	return DB
}

func CloseDB() {
	DB.Close()
}
