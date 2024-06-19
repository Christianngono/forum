package forum

import (
	"text/template"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int       `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}


func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(fmt.Sprintf("/home/christian/forum/forum/templates/%s", tmpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	} 
}

func getSessionStore() *sessions.CookieStore {
	// Charger la clé secrète à partir des variables d'environnement
	secretKey := os.Getenv("SESSION_SECRET_KEY")
	if secretKey == "" {
		secretKey = "defaultSecretKey"
	}
	return sessions.NewCookieStore([]byte(secretKey))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/christian/forum/forum/templates/index.html")
	
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	store := getSessionStore()
	session, _ := store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.ServeFile(w, r, "/home/christian/forum/forum/templates/register.html")
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

	stmt, err := DB.Prepare(`INSERT INTO users (email, username, password) VALUES (?, ?, ?)`)
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
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour gérer la connexion
	if r.Method == http.MethodPost {
		http.ServeFile(w, r, "/home/christian/forum/forum/templates/login.html")
		return
	}

	// Lire les données du formulaire
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user User
	err := DB.QueryRow("SELECT id, email, username, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return	
	}

	// Comparer les mots de passe
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Créer une sessionID à l' UUID
	sessionID, _ := uuid.NewRandom()

	store := getSessionStore()
	session, _ := store.Get(r, "session")

	session.Values["user_id"] = user.ID
	session.Values["session_id"] = sessionID.String()
	session.Save(r, w)

	http.Redirect(w, r, "/get-posts", http.StatusSeeOther)
}
