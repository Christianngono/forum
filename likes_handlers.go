package forum

import (
	"net/http"
	"strconv"
	"time"
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

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		var existingLike int
		err = DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingLike)
		if err != nil {
			http.Error(w, "Error checking like", http.StatusInternalServerError)
			return
		}
		if existingLike > 0 {
			http.Error(w, "Already liked", http.StatusBadRequest)
			return
		}

		_, err = DB.Exec("INSERT INTO likes (user_id, post_id, created_at) VALUES (?,?,?)", userID, postID, time.Now())
		if err != nil {
			http.Error(w, "Error liking post", http.StatusInternalServerError)
			return
		}

		_, err = DB.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Error updating post likes", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
	}
}
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
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

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		var existingDislike int
		err = DB.QueryRow("SELECT COUNT(*) FROM dislikes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingDislike)
		if err != nil {
			http.Error(w, "Error checking dislike", http.StatusInternalServerError)
			return
		}
		if existingDislike > 0 {
			http.Error(w, "Already disliked", http.StatusBadRequest)
			return
		}

		_, err = DB.Exec("INSERT INTO dislikes (user_id, post_id, created_at) VALUES (?,?,?)", userID, postID, time.Now())
		if err != nil {
			http.Error(w, "Error disliking post", http.StatusInternalServerError)
			return
		}

		_, err = DB.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, "Error updating post dislikes", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
	}
}
