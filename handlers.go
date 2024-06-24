package forum

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_KEY")))

var templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	posts := []Post{}
	// Récupérer posts de database
	rows, err := DB.Query("SELECT * FROM posts")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.Dislikes)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Error getting session", http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		templates.ExecuteTemplate(w, "index.html", posts)
		return
	}
	user := User{}
	err = DB.QueryRow("SELECT * FROM users WHERE id =?", userID).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	templates.ExecuteTemplate(w, "index.html", posts)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Error getting session", http.StatusInternalServerError)
		return
	}

	delete(session.Values, "user_id")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		_, err = DB.Exec("INSERT INTO users (email, username, password) VALUES (?,?,?)", email, username, hashedPassword)
		if err != nil {
			http.Error(w, "Error inserting user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := templates.ExecuteTemplate(w, "register.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user User
		err := DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			fmt.Println(err)
			return
		}
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, "Error getting session", http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = user.ID
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
