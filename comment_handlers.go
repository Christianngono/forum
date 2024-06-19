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
	if err!= nil {
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

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        renderTemplate(w, "edit-comment.html", nil)
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
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	log.Println("Received post_id:", postID)
	// Code pour récupérer les commentaires
	rows, err := DB.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", postID)
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
	json.NewEncoder(w).Encode(comments)
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
	if err!= nil {
		if err == sql.ErrNoRows {
            http.Error(w, "Comment not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
	}
	json.NewEncoder(w).Encode(comment)
	renderTemplate(w, "comment.html", comment)		   
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour supprimer un commentaire
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }
    log.Println("Received comment_id:", id)

    stmt, err := DB.Prepare("DELETE FROM comments WHERE id =?")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/get-post?id="+strconv.Itoa(id), http.StatusSeeOther)
}
