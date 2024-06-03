package forum

import (
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

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	// Log the form values for debugging
	log.Println("Form Values:", r.Form)

	

	// Get post _id form the form values
	postIDStr := r.FormValue("post_id")
	log.Println("Received post_id:", postIDStr)
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		log.Println("Invalid post_id:", postIDStr)
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

    // Get user_id form form values
	userIDStr := r.FormValue("user_id")
	log.Println("Received user_id:", userIDStr)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		log.Println("Invalid user_id:", userIDStr)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Valider si le post_id existe
	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Println("Post not found:", postID)
		http.Error(w, "Post not found", http.StatusNotFound)
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

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour récupérer un commentaire
    commentIDStr := r.URL.Query().Get("id")
    log.Println("Received comment_id:", commentIDStr)
    commentID, err := strconv.Atoi(commentIDStr)
    if err != nil || commentID <= 0 {
        log.Println("Invalid comment_id:", commentIDStr)
        http.Error(w, "Invalid comment_id", http.StatusBadRequest)
        return
    }

    var comment Comment
    err = DB.QueryRow("SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = ?", commentID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err != nil {
        log.Println("Database error:", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
