package forum

import (
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
	Likes    int`json:"likes"`
	Dislikes int `json:"dislikes"`
}



func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := getSessionStore().Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "User not logged in", http.StatusUnauthorized)
        return
    }

	if r.Method == http.MethodPost {
		http.ServeFile(w, r, "/home/christian/forum/forum/templates/create_post.html")
		return
	}

	var post Post
	post.UserID = userID
	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")

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

	stmt, err := DB.Prepare("INSERT INTO posts (id, user_id, title, content, created_at) VALUES (?, ?, ?, ?)")
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
	renderTemplate(w, "posts.html", posts)
}
func EditPostHandler(w http.ResponseWriter, r *http.Request){
	session, _ := getSessionStore().Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "User not logged in", http.StatusUnauthorized)
        return
    }

    if r.Method == http.MethodPost {
        http.ServeFile(w, r, "/home/christian/forum/forum/templates/edit_post.html")
        return
    }

    postIDStr := r.URL.Query().Get("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil || postID <= 0 {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post Post
	err = DB.QueryRow("SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts WHERE id =?", postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err!= nil {
		http.Error(w, "Post not found", http.StatusNotFound)
        return
	}
	if post.UserID != userID {
		http.Error(w, "Post not found", http.StatusNotFound)
        return
	}
	renderTemplate(w, "edit_post.html", post)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
	}

	var post Post
	err = DB.QueryRow(`SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts WHERE id = ?`, postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err!= nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return	
    }
	renderTemplate(w, "post.html", post)	
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := getSessionStore().Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "User not logged in", http.StatusUnauthorized)
        return
    }

    postIDStr := r.URL.Query().Get("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil || postID <= 0 {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post Post
	err = DB.QueryRow("SELECT id, user_id, title, content, created_at, likes, dislikes FROM posts WHERE id =?", postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err!= nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }
	if post.UserID != userID {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }
	stmt, err := DB.Prepare("DELETE FROM posts WHERE id =?")
	if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	_, err = stmt.Exec(postID)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	http.Redirect(w, r, "/get-posts", http.StatusSeeOther)
}
	

