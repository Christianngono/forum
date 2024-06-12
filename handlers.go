package forum

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Posts    []Post `json:"posts"`
	Comments []Comment `json:"comments"`
}

var templates = template.Must(template.ParseFiles("../templates/index.html",
	"../templates/register.html",
	"../templates/login.html",
	"../templates/create-post.html",
	"../templates/posts.html",
	"../templates/comments.html",
	"../templates/create-comment.html",
	"../templates/post.html",
	"../templates/comment.html",
	"../templates/likes.html",
	"../templates/dislikes.html",
))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Rendered template:", tmpl)
}

func getSessionStore() *sessions.CookieStore {
	// Charger la clé secrète à partir des variables d'environnement
	secretKey := os.Getenv("SESSION_SECRET_KEY")
	if secretKey == "" {
		// Si la clé n'est pas définie, retourner une erreur (ou utiliser une clé par défaut pour le développement)
		// A remplacer par une vraie clé en production
		println("Clé pas défini")
	}
	return sessions.NewCookieStore([]byte(secretKey))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	store := getSessionStore()
	session, _ := store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderTemplate(w, "register.html", nil)
		return
	}
	var user User
	// Décoder le formulaire envoyé
	user.Email = r.FormValue("email")
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	// Hacher le mot de passe avant de le stocker dans la base de données
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil { 
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	stmt, err := DB.Prepare("INSERT INTO users (email, username, password,) VALUES (?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour gérer la connexion
	if r.Method != http.MethodPost {
		renderTemplate(w, "login.html", nil)
		return
	}

	// Lire les données du formulaire
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user User
	err := DB.QueryRow("SELECT id, email, username, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Comparer les mots de passe
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)

	}

	// Créer une sessionID à l' UUID
	sessionID := uuid.New().String()

	// Initialiser le store de session
	store := getSessionStore()
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Save(r, w)

	// Stocker l'ID de session dans les valeurs de session
	session.Values["post_id"] = user.ID
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	session.Values["session_id"] = sessionID
	session.Save(r, w)

	/*w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User logged in successfully"})*/

	http.Redirect(w, r, "/get-posts", http.StatusSeeOther)
}
