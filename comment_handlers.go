package forum

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	CreatedAt time.Time
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Créer un commentaire
	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "session")
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		content := r.FormValue("content")

		_, err = DB.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
		if err != nil {
			http.Error(w, "Error creating comment", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/get-post?id=%d", postID), http.StatusSeeOther)
		return
	}
	if err := templates.ExecuteTemplate(w, "create-comment.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		commentID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		content := r.FormValue("content")

		_, err = DB.Exec("UPDATE comments SET content = ? WHERE id = ?", content, commentID)
		if err != nil {
			http.Error(w, "Error updating comment", http.StatusInternalServerError)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/get-post?id=%d", postID), http.StatusSeeOther)
		return
	}
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var comment Comment
	err = DB.QueryRow("SELECT id, post_id, content, FROM comments WHERE id = ?", commentID).Scan(&comment.ID, &comment.PostID, &comment.Content)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	if err := templates.ExecuteTemplate(w, "edit-comment.html", comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	rows, err := DB.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}
	if err := templates.ExecuteTemplate(w, "comments.html", comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {

	commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var comment Comment
	err = DB.QueryRow("SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = ?", commentID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	if err := templates.ExecuteTemplate(w, "comment.html", comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func FilterCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Filtrer les commentaires d'un post
	comments := []Comment{}
	rows, err := DB.Query("SELECT id, post_id, user_id, content, created_at FROM comments")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}
	if err := templates.ExecuteTemplate(w, "comments.html", comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Mettre à jour un commentaire
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	_, err = DB.Exec("UPDATE comments SET content = ? WHERE id = ?", content, commentID)
	if err != nil {
		http.Error(w, "Error updating comment", http.StatusInternalServerError)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/get-post?id=%d", postID), http.StatusSeeOther)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	//Supprimer un commentaire
	commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	_, err = DB.Exec("DELETE FROM comments WHERE id = ?", commentID)
	if err != nil {
		http.Error(w, "Error deleting comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/get-post?id=%d", postID), http.StatusSeeOther)
}
