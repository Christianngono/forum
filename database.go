package forum

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
		return
	}
	log.Println("Connected to database")

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
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
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
	createCategoriesTable := `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`
	createPostCategoriesTable := `CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER,
		category_id INTEGER,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`
	createLikesTable := `CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		comment_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id)
	);`
	createDislikesTable := `CREATE TABLE IF NOT EXISTS dislikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		comment_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id)
	);`

	statements := []string{
		createUsersTable,
		createPostsTable,
		createCommentsTable,
		createCategoriesTable,
		createPostCategoriesTable,
		createLikesTable,
		createDislikesTable,
	}

	for _, statement := range statements {
		_, err = DB.Exec(statement)
		if err != nil {
			log.Fatalf("Error creating tables: %v\n", err)
		}
		log.Printf("Created table: %s", statement)
	}
}
