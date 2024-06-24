package forum

import (
	"fmt"
	"net/http"
	"strconv"
)

type Like struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type Dislike struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		query := `UPDATE posts SET likes = likes + 1 WHERE id = ?`
		_, err = DB.Exec(query, postID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = DB.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
		if err != nil {
			http.Error(w, "Error creating like", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Post liked successfully")

		http.Redirect(w, r, fmt.Sprintf("/get-post?id=%d", postID), http.StatusSeeOther)
		return
	}
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		query := `UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?`
		_, err = DB.Exec(query, postID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = DB.Exec("INSERT INTO dislikes (user_id, post_id) VALUES (?, ?)", userID, postID)
		if err != nil {
			http.Error(w, "Error creating dislike", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Post disliked successfully")
		http.Redirect(w, r, fmt.Sprintf("/get-post?id=%d", postID), http.StatusSeeOther)
		return
	}
}
