package forum

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    
	Email    string 
	Username string 
	Password string 
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_KEY")))

var templates = template.Must(template.ParseGlob(filepath.Join("..", "templates", "*.html")))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
    if session.Values["user_id"] != nil {
        http.Redirect(w, r, "/posts", http.StatusSeeOther)
        return
    }
    templates.ExecuteTemplate(w, "index.html", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Error getting session", http.StatusInternalServerError)
		return
	}

	delete(session.Values, "user_id")
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, "Error getting session", http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = user.ID
		if err := session.Save(r, w); err!= nil {
            http.Error(w, "Error saving session", http.StatusInternalServerError)
            return
        }
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
