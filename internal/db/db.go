package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/storage/sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

func InitDb() (*sql.DB, *sqlite3.Storage) {
	// Init SQLite3 database
	// "file:memdb1?mode=memory&cache=shared" sqlite in-memory
	db, err := sql.Open("sqlite3", "./app_data.db")
	if err != nil {
		log.Fatalf("🔥 failed to connect to the database: %s", err)
	}

	db.SetConnMaxIdleTime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	stmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatalf(
			"🔥 could not create table 'users' in database: %s", err.Error(),
		)
	}

	// Storage package can create this table for you at init time
	// but for the purpose of this example I created it manually
	// expanding its structure with an "u" column to better query
	// all user-related sessions.
	stmt = `CREATE TABLE IF NOT EXISTS sessions (
		k  VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
		v  BLOB NOT NULL,
		e  BIGINT NOT NULL DEFAULT '0',
		u  TEXT);`
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatalf(
			"🔥 could not create table 'sessions' in database: %s", err.Error(),
		)
	}

	// Init sessions store
	storage := sqlite3.New(sqlite3.Config{
		Database:        "./app_data.db",
		Table:           "sessions",
		Reset:           false,
		GCInterval:      10 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 1 * time.Second,
	})

	log.Println("🚀 connected successfully to the database")

	return db, storage
}

/*
sqlite in-memory:
https://stackoverflow.com/questions/77134000/intermittent-table-missing-error-in-sqlite-memory-database
*/
