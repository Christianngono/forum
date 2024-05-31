package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour la création de commentaire
	if r.Method != http.MethodPost {
		renderTemplate(w, "create-comment.html", nil)
		return
	}

	// Log pour débogage
	log.Println("Form Values:", r.Form)

	postIDStr := r.FormValue("post_id")
	log.Println("Received post_id:", postIDStr)

	// Lire les données
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		log.Println("Invalid post_id:", postIDStr)
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}
	// Valider si le post_id existe
	var exists sql.NullBool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists.Valid || !exists.Bool {
		log.Println("Post not found:", postID)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil || userID <= 0 {
		log.Println("Invalid user_id:", r.FormValue("user_id"))
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	comment := Comment{
		PostID:    postID,
		UserID:    userID,
		Content:   r.FormValue("content"),
		CreatedAt: time.Now(),
	}

	stmt, err := DB.Prepare("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Prepare statement error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)
	if err != nil {
		log.Println("Exec statement error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	renderTemplate(w, "comments.html", nil)
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour récupérer les commentaires
	rows, err := DB.Query("SELECT id, post_id, user_id, content, created_at FROM comments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
