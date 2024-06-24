package forum

import (
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "session")
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		_, err := DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?,?,?)", userID, title, content)
		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	if err := templates.ExecuteTemplate(w, "create_post.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		_, err = DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", title, content, postID)
		if err != nil {
			http.Error(w, "Error updating post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	var post Post
	err = DB.QueryRow("SELECT id, title, content FROM posts WHERE id = ?", postID).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if err := templates.ExecuteTemplate(w, "edit_post.html", post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	_, err = DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", title, content, postID)
	if err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les postes dans database
	posts := []Post{}
	rows, err := DB.Query("SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts")
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes); err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := templates.ExecuteTemplate(w, "posts.html", posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post Post
	err = DB.QueryRow("SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts WHERE id = ?", postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
	}
	if err := templates.ExecuteTemplate(w, "post.html", post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func FilterPostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les postes dans database
	posts := []Post{}
	rows, err := DB.Query("SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes); err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}
	if err := templates.ExecuteTemplate(w, "posts.html", posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	_, err = DB.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
