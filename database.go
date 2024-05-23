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
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL,
        username TEXT NOT NULL,
        password TEXT NOT NULL
	);`
	createPostsTable := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)

	);`
	createCommentsTable := `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	_, err := db.Exec(createUsersTable)
	if err!= nil {
        log.Fatalf("Error creating users table: %v\n", err)
    }

	_, err = db.Exec(createPostsTable)
	if err!= nil {
        log.Fatalf("Error creating posts table: %v\n", err)
    }

	_, err = db.Exec(createCommentsTable)
	if err!= nil {
        log.Fatalf("Error creating comments table: %v\n", err)
    }
}

func Close() {
	db.Close()
}
