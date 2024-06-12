package forum

import (
	"database/sql"
    "strconv"
	"log"
	"net/http"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes int           `json:"likes"`
	Dislikes int        `json:"dislikes"`
}
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}


func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderTemplate(w, "create-template", r)
		return
	}
	session, _ := getSessionStore().Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "User not logged in", http.StatusUnauthorized)
        return
    }

	var post Post
	post.UserID = userID
	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	post.CreatedAt = time.Now()

	// Enregistrer les valeurs du formulaire pour le débogage
	log.Println("Form Values:", r.Form)

	// Vérifier si UserID est valide
	if post.UserID <= 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	log.Println("Parsed id:", post.ID)
	log.Println("Parsed user_id:", post.UserID)
	log.Println("Parsed title:", post.Title)
	log.Println("Parsed content:", post.Content)
	log.Println("Parsed created_at:", post.CreatedAt)

	stmt, err := DB.Prepare("INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.UserID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/get-posts", http.StatusSeeOther)   
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
		
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "posts.html", posts)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
	}
	log.Println("Received post_id:", id)

	var post Post
	err = DB.QueryRow("SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts WHERE id =?", id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err!= nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return	
	}
	renderTemplate(w, "post.html", post)
}
