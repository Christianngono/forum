package forum

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
}

func CreateTables() {
	// Code pour créer les tables dans la base de données
}

func Close() {
	db.Close()
}
