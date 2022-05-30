package database

import (
	"database/sql"
	"fmt"
	"log"
)

var (
	dbfile = "recihub.db"
)

const createProfileTable = `
	CREATE TABLE IF NOT EXISTS profile (
  	id TEXT NOT NULL PRIMARY KEY,
  	username TEXT NOT NULL UNIQUE
  );`

const createRecipeTable = `
	CREATE TABLE IF NOT EXISTS recipe (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		imagename TEXT default "",
		CHECK (name <> '' AND username <> '')
	);`

const createIngredientTable = `
	CREATE TABLE IF NOT EXISTS ingredient (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		amount TEXT NOT NULL,
		unit TEXT NOT NULL,
		recipeid INTEGER NOT NULL,
		FOREIGN KEY(recipeid) REFERENCES recipe(id) ON DELETE CASCADE
	);`

const createStepTable = `
	CREATE TABLE IF NOT EXISTS step (
		stepnumber INTEGER NOT NULL,
		description TEXT NOT NULL,
		recipeid INTEGER NOT NULL,
		PRIMARY KEY(stepnumber, recipeid),
		FOREIGN KEY(recipeid) REFERENCES recipe(id) ON DELETE CASCADE
	);`

func Open() (*sql.DB, error) {
	connName := fmt.Sprintf("%v?_foreign_keys=on", dbfile)

	db, err := sql.Open("sqlite3", connName)
	if err != nil {
		return nil, err
	}

	createSQLiteTables(db)

	return db, nil
}

func createSQLiteTables(conn *sql.DB) {
	if _, err := conn.Exec(createProfileTable); err != nil {
		log.Fatalf("failed to create PROFILE table: %s", err)
	}

	if _, err := conn.Exec(createRecipeTable); err != nil {
		log.Fatalf("failed to create RECIPE table: %s", err)
	}

	if _, err := conn.Exec(createIngredientTable); err != nil {
		log.Fatalf("failed to create INGREDIENT table: %s", err)
	}

	if _, err := conn.Exec(createStepTable); err != nil {
		log.Fatalf("failed to create STEP table: %s", err)
	}
}
