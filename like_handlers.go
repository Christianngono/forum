package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Like struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type Dislike struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

// Initialisation de la connexion à la base de données
func init() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		log.Fatal(err)
	}

	// Vérifiez la connexion
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// Handler pour gérer les likes sur les posts
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var like Like
	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", like.UserID, like.PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render template for likes
	renderTemplate(w, "likes.html", nil) // Change "likes.html" to your actual template name
}

// Handler pour gérer les dislikes sur les posts
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var dislike Dislike
	err := json.NewDecoder(r.Body).Decode(&dislike)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_, err = db.Exec("INSERT INTO dislikes (user_id, post_id) VALUES (?, ?)", dislike.UserID, dislike.PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render template for dislikes
	renderTemplate(w, "dislikes.html", nil)
}
