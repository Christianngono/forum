package forum

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour la création de commentaire
	if r.Method != http.MethodPost {
		renderTemplate(w, "create-comment.html", nil)
		return
	}

	// Analyser les données du formulaire
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enregistrez les valeurs du formulaire pour le débogage
	log.Println("Form Values:", r.Form)

	session, _ := getSessionStore().Get(r, "session")
	userID, ok := session.Values["user_id"].(int)

	if !ok {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}
	log.Println("Parsed user_id:", userID)

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}
	log.Println("Parsed post_id:", postID)

	var comment Comment
	comment.UserID = userID
	comment.PostID = postID
	comment.Content = r.FormValue("content")
	comment.CreatedAt = time.Now()

	stmt, err := DB.Prepare("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/get-post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	log.Println("Received post_id:", postID)
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
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "comments.html", comments)
}

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour récupérer un commentaire
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	log.Println("Received comment_id:", id)

	var comment Comment
	err = DB.QueryRow("SELECT id, post_id, user_id, content, created_at FROM comments WHERE id =?", id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	renderTemplate(w, "comment.html", comment)
}
