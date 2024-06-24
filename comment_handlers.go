package forum

import (
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
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, "Error getting session", http.StatusUnauthorized)
			return
		}

		userID, ok := session.Values["user_id"].(int)
		if!ok {
            http.Error(w, "User not logged in", http.StatusUnauthorized)
            return
        }

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		content := r.FormValue("content")
		_, err = DB.Exec("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", postID, userID, content, time.Now())
		if err != nil {
			http.Error(w, "Error creating comment", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
	}
}

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, "Error getting session", http.StatusInternalServerError)
			return
		}
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")

		var existingComment Comment
		err = DB.QueryRow("SELECT user_id, post_id FROM comments WHERE id = ?", commentID).Scan(&existingComment.UserID, &existingComment.PostID)
		if err != nil {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		if existingComment.UserID != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = DB.Exec("UPDATE comments SET content = ? WHERE id = ?", content, commentID)
		if err != nil {
			http.Error(w, "Error updating comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post?id="+strconv.Itoa(existingComment.PostID), http.StatusSeeOther)
	}		
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
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
			http.Error(w, "Error scanning comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating comments", http.StatusInternalServerError)
		return
	}
	if err := templates.ExecuteTemplate(w, "comments.html", comments); err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }	
}

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
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
			http.Error(w, "Error scanning comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating comments", http.StatusInternalServerError)
		return
	}
	if err := templates.ExecuteTemplate(w, "comments.html", comments); err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }	
}

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Mettre à jour un commentaire
	if r.Method == http.MethodPost {
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, "Error getting session", http.StatusInternalServerError)
			return
		}
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")

		var existingComment Comment
		err = DB.QueryRow("SELECT user_id, post_id FROM comments WHERE id = ?", commentID).Scan(&existingComment.UserID, &existingComment.PostID)
		if err != nil {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		if existingComment.UserID != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = DB.Exec("UPDATE comments SET content = ?, created_at = ? WHERE id = ?", content, time.Now(), commentID)
		if err != nil {
			http.Error(w, "Error updating comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post?id="+strconv.Itoa(existingComment.PostID), http.StatusSeeOther)
	}
	
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	//Supprimer un commentaire
	if r.Method == http.MethodPost {
		sesion, err := store.Get(r, "session")
		if err != nil {
            http.Error(w, "Error getting session", http.StatusUnauthorized)
            return
        }
		userID, ok := sesion.Values["user_id"].(int)
		if!ok {
            http.Error(w, "User not logged in", http.StatusUnauthorized)
            return
        }
		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
            http.Error(w, "Invalid comment ID", http.StatusBadRequest)
            return
        }
		var existingComment Comment 
		err = DB.QueryRow("SELECT user_id, post_id FROM comments WHERE id = ?", commentID).Scan(&existingComment.UserID, &existingComment.PostID)
		if err != nil {
            http.Error(w, "Comment not found", http.StatusNotFound)
            return
        }
		if existingComment.UserID != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
		}
		_, err = DB.Exec("DELETE FROM comments WHERE id =?", commentID)
		if err != nil {
            http.Error(w, "Error deleting comment", http.StatusInternalServerError)
            return
        }
		http.Redirect(w, r, "/post?id="+strconv.Itoa(existingComment.PostID), http.StatusSeeOther)
	}
}
